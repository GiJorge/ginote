<template>
  <div class="animate-in fade-in duration-500">
    <h1 class="md:text-6xl text-3xl font-black sm:mb-10 text-slate-900 dark:text-white leading-tight">
      {{ title }}
    </h1>

    <article class="prose dark:prose-invert max-w-none">
      <div class="text-lg" v-html="previewHtml"></div>
    </article>

    <div v-if="amharicContent" 
         class="sm:mt-16 sm:p-10 p-2 bg-slate-50 dark:bg-slate-900/50 sm:rounded-[40px] border border-slate-100 dark:border-slate-800">
      <h4 class="text-indigo-500 font-bold text-xs uppercase tracking-widest mb-6">Amharic Version</h4>
      <article class="prose dark:prose-invert max-w-none">
        <div class="text-lg" v-html="previewHtmlAmh"></div>
      </article>
    </div>
  </div>
</template>

<script setup>
defineProps(['title', 'previewHtml', 'amharicContent', 'previewHtmlAmh']);
</script>