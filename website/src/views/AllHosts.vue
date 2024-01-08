<template>
  <div>
    <ul>
      <li v-for="(host, index) in hosts" :key="index">
        {{ host.host }}
      </li>
    </ul>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      hosts: []
    };
  },
  async created() {
    const response = await axios.get('http://localhost:8081/api/allhosts');
    this.hosts = response.data;
  }
}
</script>