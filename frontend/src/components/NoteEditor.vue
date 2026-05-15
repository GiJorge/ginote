<template>
  <div class="grid grid-cols-1 gap-6">
    <input v-model="localNote.title" placeholder="Title" class="w-full text-4xl font-bold outline-none border-none bg-transparent" />
    
    <CategorySelect v-model="localNote.category" :allCategories="categories" />

    <span class="text-xl font-bold">English Section</span>
    <textarea v-model="localNote.content" placeholder="English content..." 
              class="w-full bg-transparent text-lg leading-relaxed border-2 border-slate-200 dark:border-slate-800 rounded p-4 h-64 outline-none focus:border-indigo-500 transition-colors"></textarea>

    <h1 class="text-xl font-bold">Amharic Section</h1>
    <textarea v-model="localNote.amharic_content" placeholder="Amharic content..." 
              class="w-full bg-transparent text-lg leading-relaxed italic border-2 border-slate-200 dark:border-slate-800 rounded p-4 h-64 outline-none focus:border-indigo-500 transition-colors"></textarea>
  </div>
</template>

<script setup>
import CategorySelect from './CategorySelect.vue';
const props = defineProps(['modelValue', 'categories']);
const localNote = props.modelValue; // Reactive reference back to App.vue
</script>