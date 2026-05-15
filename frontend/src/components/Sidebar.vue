<template>
  <aside :class="['fixed inset-y-0 left-0 z-50 w-72 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 transform transition-transform duration-300 ease-in-out md:relative md:translate-x-0 flex flex-col', isOpen ? 'translate-x-0' : '-translate-x-full']">
    
    <div class="p-6 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center shrink-0">
      <h2 class="font-bold text-xl text-indigo-600">NoteHub</h2>
      <div class="flex gap-2">
        <button @click="$emit('toggle-theme')" class="p-2 rounded-full bg-slate-100 dark:bg-slate-800">{{ isDark ? '🌙' : '☀️' }}</button>
        <button @click="$emit('add-note')" class="bg-indigo-600 text-white p-2 rounded-lg">＋</button>
      </div>
    </div>

    <div class="p-4 space-y-2 shrink-0">
      <input v-model="search" type="text" placeholder="Search..." class="w-full p-2 bg-slate-100 dark:bg-slate-800 rounded-lg outline-none" />
      <select v-model="filter" class="w-full p-2 bg-slate-100 dark:bg-slate-800 rounded-lg text-sm outline-none">
        <option value="">All Categories</option>
        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
      </select>
    </div>

    <div class="overflow-y-auto flex-1 p-3 space-y-2 custom-scrollbar ">
      <div v-for="n in filteredNotes" :key="n.id" @click="$emit('select-note', n)"
           :class="['p-4 rounded-xl cursor-pointer transition-colors', activeId === n.id ? 'bg-indigo-600 text-white shadow-lg' : 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-300']">
        <h3 class="font-semibold truncate">{{ n.title || 'Untitled' }}</h3>
        <div class="flex justify-between items-center mt-1">
          <span class="text-[10px] uppercase font-bold opacity-60">{{ n.category || 'General' }}</span>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue';
const props = defineProps(['isOpen', 'isDark', 'notes', 'categories', 'activeId']);
const emit = defineEmits(['toggle-theme', 'add-note', 'select-note']);

const search = ref("");
const filter = ref("");

const filteredNotes = computed(() => {
  return props.notes.filter(n => {
    const term = search.value.toLowerCase();
    const matchesSearch = n.title.toLowerCase().includes(term) || n.content.toLowerCase().includes(term);
    const matchesCategory = filter.value === "" || n.category === filter.value;
    return matchesSearch && matchesCategory;
  });
});
</script>