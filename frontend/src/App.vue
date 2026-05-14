<template>
 <!-- Root div with 'dark' class binding -->
  <div :class="{ 'dark': isDark }">
    <div class="flex h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 transition-colors duration-300">

    <!-- Sidebar -->
  <aside class="w-80 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col">
     <div class="p-6 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center">
              <h2 class="font-bold text-xl text-indigo-600">NoteHub</h2>
      <div class="flex gap-2">
             <button @click="toggleTheme" class="p-2 rounded-full bg-slate-100 dark:bg-slate-800 transition-colors">
    {{ isDark ? '🌙' : '☀️' }}
  </button>

            <button @click="resetForm" class="bg-indigo-600 text-white p-2 rounded-lg">＋</button>
          </div>   
         </div>



   <!-- Search & Filter Area -->
        <div class="p-4 space-y-2">
          <input v-model="searchQuery" type="text" placeholder="Search..."
                 class="w-full p-2 bg-slate-100 dark:bg-slate-800 rounded-lg outline-none" />

          <select v-model="filterCategory" class="w-full p-2 bg-slate-100 dark:bg-slate-800 rounded-lg text-sm outline-none">
            <option value="">All Categories</option>
            <option v-for="cat in uniqueCategories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>

        <!-- Sidebar List: Uses FILTERED notes -->
        <div class="overflow-y-auto flex-1 p-3 space-y-2">
          <div v-for="n in filteredNotes" :key="n.id" @click="selectNote(n)"
               :class="['p-4 rounded-xl cursor-pointer', currentNoteId === n.id ? 'bg-indigo-600 text-white' : 'hover:bg-slate-100 dark:hover:bg-slate-800']">
            <h3 class="font-semibold truncate">{{ n.title || 'Untitled' }}</h3>
            <span class="text-[10px] uppercase font-bold opacity-60">{{ n.category || 'General' }}</span>
          </div>
        </div>



    </aside>

    <!-- Main Content -->
   <main class="flex-1 flex flex-col bg-white dark:bg-slate-950 overflow-hidden">
  
  <!-- HEADER: Changes based on Mode -->
  <header class="h-16 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center px-8 shrink-0">
    <!-- Show Edit/Preview Toggles ONLY when in Edit/Preview mode -->
    <div v-if="viewMode !== 'read'" class="flex bg-slate-100 dark:bg-slate-800 p-1 rounded-lg">
      <button @click="viewMode = 'edit'" :class="['px-4 py-1 text-sm rounded-md', viewMode === 'edit' ? 'bg-white dark:bg-slate-700 shadow font-bold' : '']">Write</button>
      <button @click="viewMode = 'preview'" :class="['px-4 py-1 text-sm rounded-md', viewMode === 'preview' ? 'bg-white dark:bg-slate-700 shadow font-bold' : '']">Preview</button>
    </div>
    
    <!-- Show a simple 'Back' or 'Close' if in Read mode -->
    <div v-else>
      <span class="text-xs font-bold text-slate-400 uppercase tracking-widest">Reading Mode</span>
    </div>

    <div class="flex gap-4">
      <!-- Edit Button: Only shows in Read Mode -->
      <button v-if="viewMode === 'read'" @click="startEditing" class="px-6 py-2 border border-slate-200 dark:border-slate-700 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-800 transition-all">
        Edit Note
      </button>

         <button v-if="viewMode === 'preview'" @click="translateNote" :disabled="isTranslating" class="bg-purple-600 text-white px-4 py-2 rounded">
  {{ isTranslating ? 'Translating...' : 'Translate to Amharic 🇪🇹' }}
</button>





          <button v-if="currentNoteId" @click="deleteNote" class="bg-red-500 text-white px-4 py-2 rounded">Delete</button>
      
      <!-- Save Button: Shows in Edit/Preview -->
      <button v-if="viewMode !== 'read'" @click="saveNote" class="px-6 py-2 bg-indigo-600 text-white rounded-lg font-bold shadow-lg shadow-indigo-100 dark:shadow-none">
        Save Changes
      </button>
    </div>
  </header>

  <div class="flex-1 p-10 overflow-y-auto">
    <div class="max-w-4xl mx-auto">

      <!-- 1. STANDALONE READ MODE -->
      <div v-if="viewMode === 'read'" class="animate-in fade-in duration-500">
        <h1 class="text-6xl font-black mb-10 text-slate-900 dark:text-white leading-tight">{{ note.title }}</h1>
        
        <article class="prose dark:prose-invert max-w-none">
          <div class="text-lg" v-html="previewHtml"></div>
        </article>

        <div v-if="note.amharic_content" class="mt-16 p-10 bg-slate-50 dark:bg-slate-900/50 rounded-[40px] border border-slate-100 dark:border-slate-800">
          <h4 class="text-indigo-500 font-bold text-xs uppercase tracking-widest mb-6">Amharic Version</h4>
            <article class="prose dark:prose-invert max-w-none">
              <div class="text-lg" v-html="previewHtmlamh"></div>
            </article>
        </div>
      </div>



      <!-- 2. EDIT MODE (The Grid setup we built) -->

 <div v-if="viewMode === 'edit'" class="grid grid-cols-1 gap-6 text-taupe-900 dark:text-taupe-100">
        <input v-model="note.title" placeholder="Title" class="w-full text-4xl font-bold outline-none border-none bg-transparent " />
       <!--   <input v-model="note.category" placeholder="Category" class="w-full p-2 border-2 border-gray-300 dark:border-gray-900 focus:border-yellow-600 rounded" />
      --> 
     <div class="relative w-full max-w-xs">
  <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2 block">
    Category
  </label>
  
  <input 
    type="text"
    v-model="searchCatQuery"
    @focus="showDropdown = true"
    @blur="closeDropdown"
  
    placeholder="Type to search or add..."
    class="w-full bg-slate-100 dark:bg-slate-800 border-none rounded-xl px-4 py-2 text-sm focus:ring-2 focus:ring-indigo-500 outline-none transition-all"
  />

  <!-- Dropdown Menu -->
  <div v-if="showDropdown" class="absolute z-50 w-full mt-2 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-xl overflow-hidden animate-in fade-in zoom-in-95 duration-200">
    
    <!-- Existing Categories -->
    <div v-if="filteredCategories.length > 0">
      <button 
        v-for="cat in filteredCategories" 
        :key="cat"
        @click="setCategory(cat)"
        class="w-full text-left px-4 py-3 text-sm hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors border-b border-slate-100 dark:border-slate-700/30 last:border-0"
      >
        {{ cat }}
      </button>
    </div>

    <!-- "Add New" Option (Shows if the exact query doesn't exist) -->
    <button 
      v-if="searchCatQuery && !uniqueCategories.includes(searchCatQuery)"
   
      @mousedown="setCategory(searchCatQuery)"
      class="w-full text-left px-4 py-3 text-sm text-indigo-500 font-bold hover:bg-indigo-50 dark:hover:bg-indigo-500/10 transition-colors"
    >
      + Add new: "{{ searchCatQuery }}"
    </button>
    
    <div v-if="!searchCatQuery && filteredCategories.length === 0" class="px-4 py-3 text-xs text-slate-400 italic">
      Start typing to create a category...
    </div>
  </div>
</div> 
      
      
      
      
      
          <!-- English Section -->
        <span class="text-xl font-bold">English Section</span>
<div class="auto-grid-area" :data-replicated-value="note.content">
  <textarea 
    v-model="note.content" 
    placeholder="English content..."
    class="bg-transparent text-lg leading-relaxed outline-none focus:ring-0 border-2 border-gray-300 dark:border-gray-900 focus:border-green-600 rounded p-2"
  ></textarea>
</div>
<div class="relative flex py-3 items-center">
    <div class="grow border-t border-gray-400"></div>
    <span class="shrink mx-4 text-gray-400">English Section End</span>
    <div class="grow border-t border-gray-400"></div>
</div>
<h1 class="text-xl font-bold">Amharic Section</h1>
<!-- Amharic Section -->
<div class="auto-grid-area" :data-replicated-value="note.amharic_content">
  <textarea 
    v-model="note.amharic_content" 
    placeholder="Amharic content..."
    class="bg-transparent text-lg leading-relaxed italic outline-none focus:ring-0 border-2 border-gray-300 dark:border-gray-900 focus:border-blue-600 rounded p-2"
  ></textarea>
</div>
<div class="relative flex py-3 items-center">
    <div class="grow border-t border-gray-400"></div>
    <span class="shrink mx-4 text-gray-400">Amharic Section End</span>
    <div class="grow border-t border-gray-400"></div>
</div>

      </div>



      <!-- 3. PREVIEW MODE (While Editing) -->
      <div v-else-if="viewMode === 'preview'" class="prose dark:prose-invert max-w-none">
        <h1 class="text-5xl font-black mb-8">{{ note.title }}</h1>
        <div class="text-lg" v-html="previewHtml"></div><br>
<div class="relative flex py-5 items-center">
    <div class="grow border-t border-gray-400"></div>
    <span class="shrink mx-4 text-gray-400">Amharic Section </span>
    <div class="grow border-t border-gray-400"></div>
</div>
   <article class="prose dark:prose-invert prose-indigo max-w-none
                              prose-table:border prose-th:bg-slate-100 dark:prose-th:bg-slate-900
                              prose-strong:text-indigo-100 dark:prose-strong:text-indigo-400">
                <div class="text-lg" v-html="previewHtmlamh"></div>
              </article>



      </div>
 


    </div>
  </div>
</main>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted ,watch, nextTick} from 'vue';
import axios from 'axios';
import { marked } from 'marked';


const isDark = ref(true)

const toggleTheme = () => {
  isDark.value = !isDark.value
  const html = document.documentElement
  if (isDark.value) {
    html.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    html.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

onMounted(() => {
  if (localStorage.getItem('theme') === 'light') {
    isDark.value = false
    document.documentElement.classList.remove('dark')
  } else {
    document.documentElement.classList.add('dark')
  }
})








//const API_URL = '/api';

const API_URL = "http://localhost:8389/api";
const notesList = ref([]);
const isEditing = ref(true);
//const viewMode = ref('edit');

const viewMode = ref('read');

const currentNoteId = ref(null);
const note = ref({ id: 0, title: "", content: "", amharic_content: "", category: "" });

//const previewHtmlamh = computed(() => marked(note.value.amharic_content || ''));
const previewHtml = computed(() => marked(note.value.content || ''));
const previewHtmlamh = computed(() => {
  return marked.parse(note.value.amharic_content || '');
});

const searchQuery = ref("");
const filterCategory = ref("");

const engArea = ref(null);
const amhArea = ref(null);



const searchCatQuery = ref('');
const showDropdown = ref(false);


// Filter categories based on what the user types
const filteredCategories = computed(() => {
  return uniqueCategories.value.filter(cat => 
    cat.toLowerCase().includes(searchCatQuery.value.toLowerCase())
  );
});

// Select or Create handler
// const setCategory = (cat) => {
//   note.value.category = cat;
//   searchCatQuery.value = cat;
//   showDropdown.value = false;
// };

const closeDropdown = () => {
  // We use a small delay so the 'click' event on the 
  // dropdown items fires before the menu is removed from the DOM
  setTimeout(() => {
    showDropdown.value = false;
  }, 200);
};

const setCategory = (cat) => {
  // 1. Update the actual note data (what gets sent to the Go backend)
  note.value.category = cat;
  
  // 2. Update the input field text
  searchCatQuery.value = cat;
  
  // 3. Close the menu
  showDropdown.value = false;
  
  console.log("Category set to:", note.value.category); // Check your console!
};


//const note = ref({ id: 0, title: "", content: "", category: "", amharic_content: "" });

// Fix for Search: Filtering based on title or content
const filteredNotes = computed(() => {
  return notesList.value.filter(n => {
    const term = searchQuery.value.toLowerCase();
    const matchesSearch = n.title.toLowerCase().includes(term) || n.content.toLowerCase().includes(term);
    const matchesCategory = filterCategory.value === "" || n.category === filterCategory.value;
    return matchesSearch && matchesCategory;
  });
});

// Dynamic Categories for Dropdown
const uniqueCategories = computed(() => {
  const cats = notesList.value.map(n => n.category).filter(Boolean);
  return [...new Set(cats)];
});
const fetchNotes = async () => {
  const res = await axios.get(`${API_URL}/notes`);
  notesList.value = res.data || [];
};

// const selectNote = (n) => {
//   note.value = { ...n };
//   currentNoteId.value = n.id;
//   isEditing.value = false;
// };

// const selectNote = (n) => {
//   note.value = { ...n };
//   currentNoteId.value = n.id;
//   viewMode.value = 'read'; // Always open in Standalone View first
//   isEditing.value = false
// };



const selectNote = (n) => {
  // If we have a draft of the same note, load the draft instead of the list version
  const draft = localStorage.getItem('ginote_draft');
  if (draft) {
    const parsedDraft = JSON.parse(draft);
    if (parsedDraft.id === n.id) {
      note.value = parsedDraft;
      currentNoteId.value = n.id;
      viewMode.value = 'read';
      return;
    }
  }

  // Otherwise, load from the list
  note.value = { ...n };
  currentNoteId.value = n.id;
  viewMode.value = 'read';
};







const startEditing = () => {
  viewMode.value = 'edit';
};


const resetForm = () => {
  note.value = { id: 0, title: "", content: "", amharic_content: "", category: "" };
  currentNoteId.value = null;
  isEditing.value = true;
   viewMode.value = 'edit';
};

// const saveNote = async () => {
//   const res = await axios.post(`${API_URL}/notes`, note.value);
//   if (!currentNoteId.value) {
//     currentNoteId.value = res.data.id;
//     note.value.id = res.data.id;
//   }
//   await fetchNotes();
//   alert("Saved!");
// };


// 3. Clear the draft ONLY after a successful save to the database
const saveNote = async () => {
  const res = await axios.post(`${API_URL}/notes`, note.value);
  try {
    // ... your existing axios/fetch POST logic ...
    if (!currentNoteId.value) {
     currentNoteId.value = res.data.id;
     note.value.id = res.data.id;
   }
    // If the backend save is successful:
    localStorage.removeItem('ginote_draft'); 
    alert("Note saved to database!");
  } catch (error) {
    console.error("Save failed, but draft is still in local storage.");
  }
   await fetchNotes();
   //alert("Saved!");
};











const deleteNote = async () => {
  if (confirm("Delete this note?")) {
    await axios.delete(`${API_URL}/notes/${currentNoteId.value}`);
    resetForm();
    fetchNotes();
  }
};

const isTranslating = ref(false);

const translateNote = async () => {
  isTranslating.value = true;
  try {
    const res = await axios.post(`${API_URL}/notes/${currentNoteId.value}/translate`);
    note.value.amharic_content = res.data.amharic_content;
  } catch (error) {
    alert("Translation failed. The free API might be busy.");
  } finally {
    isTranslating.value = false;
  }
};








// Your existing note ref
// const note = ref({
//   title: '',
//   content: '',
//   amharic_content: '',
//   category: ''
// });

// 1. On Load: Check if there is a saved draft
onMounted(() => {
  const savedDraft = localStorage.getItem('ginote_draft');
  if (savedDraft) {
    note.value = JSON.parse(savedDraft);
  }
});

// 2. Watch for changes and save to LocalStorage automatically
watch(note, (newVal) => {
  localStorage.setItem('ginote_draft', JSON.stringify(newVal));
}, { deep: true });









/*
// Sync heights when a note is loaded or translated
watch(() => note.value.content, async () => {
  await nextTick();
  if (engArea.value) adjustHeight(engArea.value);
}, { immediate: true });

watch(() => note.value.amharic_content, async () => {
  await nextTick();
  if (amhArea.value) adjustHeight(amhArea.value);
}, { immediate: true });
*/

onMounted(fetchNotes);
</script>
