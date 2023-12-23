<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import Login from './views/LoginView.vue'

import {FCMStore} from "@/stores/fcm"
import { onMounted } from 'vue';

let fcm = FCMStore()
onMounted(async()=>{
  await fcm.SetApp();
})
</script>

<template>
  <RouterView />
</template>