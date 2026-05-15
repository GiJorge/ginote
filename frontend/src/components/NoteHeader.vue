<template>
  <header class="h-16 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center px-8 shrink-0">
    <div v-if="viewMode !== 'read'" class="flex bg-slate-100 dark:bg-slate-800 p-1 rounded-lg">
      <button @click="$emit('update:viewMode', 'edit')" 
              :class="['px-4 py-1 text-sm rounded-md', viewMode === 'edit' ? 'bg-white dark:bg-slate-700 shadow font-bold' : '']">
        Write
      </button>
      <button @click="$emit('update:viewMode', 'preview')" 
              :class="['px-4 py-1 text-sm rounded-md', viewMode === 'preview' ? 'bg-white dark:bg-slate-700 shadow font-bold' : '']">
        Preview
      </button>
    </div>

    <div v-else>
      <span class="text-xs font-bold text-slate-400 uppercase tracking-widest">Reading Mode</span>
    </div>

    <div class="flex gap-4">
      <button v-if="viewMode === 'read'" @click="$emit('update:viewMode', 'edit')" 
              class="px-6 py-2 border border-slate-200 dark:border-slate-700 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-800 transition-all">
        ✏️
      </button>

      <button v-if="viewMode === 'preview'" @click="$emit('translate')" :disabled="isTranslating" 
              class="text-white px-4 py-2 rounded bg-indigo-500 disabled:opacity-50">
        {{ isTranslating ? '...' : '🇪🇹' }}
      </button>

      <button v-if="activeId" @click="$emit('delete')" class="text-white px-4 py-2 rounded bg-red-500/10 hover:bg-red-500/20">
        ❌
      </button>

      <button v-if="viewMode !== 'read'" @click="$emit('save')" 
              class="px-6 py-2 bg-indigo-600 text-white rounded-lg font-bold shadow-lg">
        💾
      </button>
    </div>
  </header>
</template>

<script setup>
defineProps(['viewMode', 'activeId', 'isTranslating']);
defineEmits(['update:viewMode', 'translate', 'delete', 'save']);
</script>