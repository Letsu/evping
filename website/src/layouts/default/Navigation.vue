<template>
  <v-navigation-drawer>
      <v-list>
      <v-list-item title="Multi Ping tool" subtitle="or something"></v-list-item>
      <v-divider></v-divider>
      <v-list-subheader>Irgendwas 1</v-list-subheader>
      <v-list-item link title="Home"></v-list-item>
      <v-list-item 
      link title="All Hosts"
      :to="{ name: 'AllHosts' }"
      ></v-list-item>
      <v-list-subheader>Hosts</v-list-subheader>
      <v-list-item 
        v-for="(host, index) in hosts"
        :key="index"
        :to="{ name: 'Host', params: { host: host.host } }"
        link>
        {{ host.host }}
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>
  
<script>
import axios from 'axios'

export default {
  name: 'Navigation',
  components: { },
  data() {
    return {
      hosts: []
    }
  },
  async created() {
    const response = await axios.get('http://localhost:8081/api/allhosts');
    this.hosts = response.data;
  }
}
</script>
  