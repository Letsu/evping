<template>
  <h1>{{ host }}</h1>
</template>

<script setup>
  import { ref, onMounted } from 'vue';
  import { useRoute } from 'vue-router';

  const route = useRoute();
  const host = ref('');

  onMounted(() => {
    host.value = route.params.host;
  });
</script>