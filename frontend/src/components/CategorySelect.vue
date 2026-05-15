<template>
  <div class="relative w-full max-w-xs">
    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2 block">Category</label>
    <input
      type="text"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      @focus="show = true"
      @blur="close"
      placeholder="Search or add..."
      class="w-full bg-slate-100 dark:bg-slate-800 border-none rounded-xl px-4 py-2 text-sm focus:ring-2 focus:ring-indigo-500 outline-none"
    />

    <div v-if="show" class="absolute z-50 w-full mt-2 bg-white dark:bg-slate-800 border border-slate-200 rounded-2xl shadow-xl overflow-hidden">
      <button v-for="cat in filtered" :key="cat" @mousedown="$emit('update:modelValue', cat); show = false"
              class="w-full text-left px-4 py-3 text-sm hover:bg-slate-50 dark:hover:bg-slate-700 border-b border-slate-100 last:border-0">
        {{ cat }}
      </button>
      <button v-if="modelValue && !allCategories.includes(modelValue)" @mousedown="show = false"
              class="w-full text-left px-4 py-3 text-sm text-indigo-500 font-bold hover:bg-indigo-50">
        + Add: "{{ modelValue }}"
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
const props = defineProps(['modelValue', 'allCategories']);
const emit = defineEmits(['update:modelValue']);
const show = ref(false);

const filtered = computed(() => 
  props.allCategories.filter(c => c.toLowerCase().includes(props.modelValue.toLowerCase()))
);
const close = () => setTimeout(() => show.value = false, 200);
</script>