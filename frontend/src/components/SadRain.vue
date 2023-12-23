<script setup lang="ts">
import { onMounted } from "vue";
// @ts-ignore
import {useReward} from "vue-rewards";

const config = {
    startVelocity:20,
    spread:1000,
    elementCount:500,
    emoji: ["😔"],
  }

onMounted(()=>{
    const { reward } = useReward('sad', 'emoji', config);

    reward();
})
</script>


<template>
   <div class="absolute" style="transform: translateX(50%); left: 50%;" id="sad"></div>
</template>

