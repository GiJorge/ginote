<template>
  <div :class="{ 'dark': isDark }">
    <div class="flex h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 transition-colors duration-300">

      <div
        v-if="isSidebarOpen"
        @click="isSidebarOpen = false"
        class="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm md:hidden"
      ></div>

      <Sidebar 
        :isOpen="isSidebarOpen" 
        :isDark="isDark" 
        :notes="notesList" 
        :categories="uniqueCategories" 
        :activeId="currentNoteId"
        @toggle-theme="toggleTheme"
        @add-note="resetForm"
        @select-note="selectNote"
      />

      <main class="flex-1 flex flex-col bg-white dark:bg-slate-950 overflow-hidden">
        
        <header class="flex items-center justify-between p-4 md:hidden border-b border-slate-200 dark:border-slate-800">
          <button @click="isSidebarOpen = true" class="p-2 text-slate-600 dark:text-slate-400">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="3" y1="12" x2="21" y2="12"></line>
              <line x1="3" y1="6" x2="21" y2="6"></line>
              <line x1="3" y1="18" x2="21" y2="18"></line>
            </svg>
          </button>
          <span class="font-bold text-indigo-600">GiNote</span>
          <div class="w-10"></div> </header>

        <NoteHeader 
          v-model:viewMode="viewMode" 
          :activeId="currentNoteId"
          :isTranslating="isTranslating"
          @save="saveNote"
          @delete="deleteNote"
          @translate="translateNote"
        />

        <div class="flex-1 sm:p-10 p-4 overflow-y-auto">
          <div class="max-w-4xl mx-auto">
            
            <NoteViewer 
              v-if="viewMode === 'read'" 
              :title="note.title"
              :previewHtml="previewHtml"
              :amharicContent="note.amharic_content"
              :previewHtmlAmh="previewHtmlamh"
            />

            <NoteEditor 
              v-else-if="viewMode === 'edit'" 
              v-model="note"
              :categories="uniqueCategories"
            />

            <NoteViewer 
              v-else-if="viewMode === 'preview'" 
              :title="note.title"
              :previewHtml="previewHtml"
              :amharicContent="note.amharic_content"
              :previewHtmlAmh="previewHtmlamh"
            />

          </div>
        </div>
      </main>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import axios from 'axios';
import { marked } from 'marked';

// Import your new components
import Sidebar from './components/Sidebar.vue';
import NoteHeader from './components/NoteHeader.vue';
import NoteEditor from './components/NoteEditor.vue';
import NoteViewer from './components/NoteViewer.vue';

// --- State Management ---
const isSidebarOpen = ref(false);
const isDark = ref(true);
const viewMode = ref('read'); // 'read', 'edit', 'preview'
const isTranslating = ref(false);

const API_URL = '/api'; // Use relative path for Docker/Debian compatibility
const notesList = ref([]);
const currentNoteId = ref(null);
const note = ref({ id: 0, title: "", content: "", amharic_content: "", category: "" });

// --- Markdown Previews ---
const previewHtml = computed(() => marked.parse(note.value.content || ''));
const previewHtmlamh = computed(() => marked.parse(note.value.amharic_content || ''));

// --- Theme Logic ---
const toggleTheme = () => {
  isDark.value = !isDark.value;
  const html = document.documentElement;
  if (isDark.value) {
    html.classList.add('dark');
    localStorage.setItem('theme', 'dark');
  } else {
    html.classList.remove('dark');
    localStorage.setItem('theme', 'light');
  }
};

// --- API Methods ---
const fetchNotes = async () => {
  try {
    const res = await axios.get(`${API_URL}/notes`);
    notesList.value = res.data || [];
  } catch (err) {
    console.error("Failed to fetch notes. Is the Go backend running on :8389?");
  }
};

const saveNote = async () => {
  try {
    const res = await axios.post(`${API_URL}/notes`, note.value);
    if (!currentNoteId.value) {
      currentNoteId.value = res.data.id;
      note.value.id = res.data.id;
    }
    localStorage.removeItem('ginote_draft'); // Clear draft on successful save
    await fetchNotes();
    alert("Saved successfully!");
  } catch (error) {
    console.error("Save failed. Check folder permissions in /opt/ginote/data.");
  }
};

const deleteNote = async () => {
  if (confirm("Delete this note?")) {
    await axios.delete(`${API_URL}/notes/${currentNoteId.value}`);
    resetForm();
    await fetchNotes();
  }
};

const translateNote = async () => {
  isTranslating.value = true;
  try {
    const res = await axios.post(`${API_URL}/notes/${currentNoteId.value}/translate`);
    note.value.amharic_content = res.data.amharic_content;
  } catch (error) {
    alert("Translation failed. Check backend logs.");
  } finally {
    isTranslating.value = false;
  }
};

// --- Form Helpers ---
const selectNote = (n) => {
  const draft = localStorage.getItem('ginote_draft');
  isSidebarOpen.value = false;

  if (draft) {
    const parsedDraft = JSON.parse(draft);
    if (parsedDraft.id === n.id) {
      note.value = parsedDraft;
      currentNoteId.value = n.id;
      viewMode.value = 'read';
      return;
    }
  }

  note.value = { ...n };
  currentNoteId.value = n.id;
  viewMode.value = 'read';
};

const resetForm = () => {
  note.value = { id: 0, title: "", content: "", amharic_content: "", category: "" };
  currentNoteId.value = null;
  viewMode.value = 'edit';
  isSidebarOpen.value = false;
};

// --- Computed Categories ---
const uniqueCategories = computed(() => {
  const cats = notesList.value.map(n => n.category).filter(Boolean);
  return [...new Set(cats)];
});

// --- Lifecycle & Watchers ---
onMounted(() => {
  // Restore Theme
  if (localStorage.getItem('theme') === 'light') {
    isDark.value = false;
    document.documentElement.classList.remove('dark');
  } else {
    document.documentElement.classList.add('dark');
  }

  // Restore Draft
  const savedDraft = localStorage.getItem('ginote_draft');
  if (savedDraft) {
    note.value = JSON.parse(savedDraft);
  }

  fetchNotes();
});

watch(note, (newVal) => {
  localStorage.setItem('ginote_draft', JSON.stringify(newVal));
}, { deep: true });

</script>