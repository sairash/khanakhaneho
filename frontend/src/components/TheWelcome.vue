<script setup lang="ts">
import { ref } from "vue";
import { AuthStore } from "@/stores/authStore";
import { DashboardStore } from "@/stores/dashboardStore";
import logo from "@/assets/khanakhaneho_white.png"
import axios, { AxiosError } from "axios";
import { showLoadingToast, closeToast, showConfirmDialog } from 'vant';
import { onMounted } from "vue";
import router from "@/router";

import SadRain from "./SadRain.vue";

import ConfettiExplosion from "vue-confetti-explosion";




const show = ref(false);
const actions = [
  { name: 'Logout', color: '#ee0a24' },
];

function logout(){
  showLoadingToast({
    duration: 0,
    message: 'Loading...',
    forbidClick: true,
  });
  authStore.deleteAuth().then(()=>{
    closeToast();
  })
}

const anchors = [
  0,
  300,
];
const height = ref(anchors[0]);
const friend_number = ref()

let authStore = AuthStore()
let dashboard = DashboardStore();

const current_message_status = ref("not_asked");


function click_button(value: string, type_of_btn: string) {
  showLoadingToast({
    duration: 0,
    message: 'Loading...',
    forbidClick: true,
  });
  axios.post(`${import.meta.env.VITE_HOST}/send/message`, {
    message: value,
    type_of_message: type_of_btn
  }, {
    headers: {
      'Authorization': `Bearer ${authStore.token}`
    }
  }).then(response => {
    if(response.data.result == 2){
      logout()
    }
    closeToast();
    get_messages()
  }).catch((error: AxiosError) => {
    console.log("in Error ", error.response?.data)
  })
}

function show_logout(){
  show.value = true;
}

function show_change_friend(){
  if(height.value == anchors[0]){
    height.value = anchors[1]
  }else{
    height.value = anchors[0]
  }
}

interface notifications {
  notifications: notification[]
}

interface notification {
  id: number,
  UserId: number,
  extra_value: number,
  message: string
}

interface result_values {
  conform_text: string
  cancel_text: string,
  message: string,
  title: string,
}

interface result {
    [key: string]: result_values;
  }

let result:result={
  friend_request: {
    conform_text: "Add Friend",
    cancel_text: "Cancle Request",
    message: "Accept Friend Request from user id ",
    title: "Friend Request!"
  }
} 

async function get_messages(){
  axios.post(`${import.meta.env.VITE_HOST}/notification`, {}, {
    headers: {
      'Authorization': `Bearer ${authStore.token}`
    }
  }).then(response => {
    if(response.data.result == 2){
      logout()
    }
    (response.data.data as notifications).notifications.forEach(element => {
      let get_res = result[element.message]
      showConfirmDialog({
        title: get_res.title,
        message: get_res.message + element.extra_value,
        cancelButtonDisabled: false,
        cancelButtonText: get_res.cancel_text,
        confirmButtonText: get_res.conform_text
      }).then(() => {
        change_friend(element.id, true, true)
      }).catch(e =>{
        change_friend(element.id, true, false)
      });
    });
  }).catch((error: AxiosError) => {
    console.log("in Error ", error.response?.data)
  })

  axios.post(`${import.meta.env.VITE_HOST}/get/message`, {}, {
    headers: {
      'Authorization': `Bearer ${authStore.token}`
    }
  }).then(response => {
    if(response.data.result == 2){
      logout()
    }
    current_message_status.value = `${response.data.data.message.type_of_message}_${response.data.data.message.message}`
    console.log(response.data.result)
  }).catch((error: AxiosError) => {
    console.log("in Error ", error.response?.data)
  })
}

onMounted(async()=>{
  let toke = await authStore.IsAuth()
  if(!toke){
    router.push({ name: 'login'})
  }

  get_messages()
})


function change_friend(id_or_number: number, accepting: boolean, applying: boolean){
  axios.post(`${import.meta.env.VITE_HOST}/change_friend`, {
    "change_to_user": id_or_number,
    "accepting": accepting,
    "applying": applying
}, {
    headers: {
      'Authorization': `Bearer ${authStore.token}`
    }
  }).then(response => {
    if(response.data.result == 2){
      logout()
    }
    location.reload()
    console.log(response.data.result)
  }).catch((error: AxiosError) => {
    console.log("in Error ", error.response?.data)
  })
}

</script>

<template>
  <div class="w-screen h-screen overflow-hidden" :class="dashboard.get_question(authStore.auth.RoleId, 'bg',current_message_status, authStore.auth.FriendRole)">
    <div class="" v-if="dashboard.get_question(authStore.auth.RoleId, 'extra',current_message_status, authStore.auth.FriendRole) == 'confetti'">
      <div class="absolute top-0 left-0">
        <ConfettiExplosion :stageHeight="1000"/>
      </div>
      <div class="absolute top-0 right-0">
        <ConfettiExplosion :stageHeight="1000"/>
      </div>
    </div>
    <div class="" v-else-if="dashboard.get_question(authStore.auth.RoleId, 'extra',current_message_status, authStore.auth.FriendRole) == 'sad'">
      <SadRain></SadRain>
    </div>
    <div class="w-full flex justify-center">
      <div class="">
        <div class="w-full max-w-sm pt-4 px-5 text-2xl font-bold flex justify-between absolute top-0">
          <svg @click="show_change_friend()" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
          </svg>
          <button @click="show_logout">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
            </svg>
          </button>

        </div>
        <div class="w-full max-w-sm h-screen flex flex-col justify-around pt-14">
          <div class="w-full">
            <img :src="logo" class="w-full px-2">
            <div class="text-center mt-2 text-lg font-semibold">{{ dashboard.get_question(authStore.auth.RoleId, "starter",
              current_message_status, authStore.auth.FriendRole) }} {{ dashboard.get_question(authStore.auth.FriendRole, "role",
    current_message_status, authStore.auth.FriendRole) }} {{ dashboard.get_question(authStore.auth.RoleId, "ender", current_message_status, authStore.auth.FriendRole) }}
            </div>
          </div>
          <div class="flex w-full justify-center">
            <img :src="dashboard.get_question(authStore.auth.RoleId, 'image', current_message_status, authStore.auth.FriendRole)" class="drop-shadow max-w-xs" alt="" srcset="">
          </div>
          <div class="w-full px-5 h-32 flex flex-col justify-center gap-3">
            <div class="" v-for="but in dashboard.get_buttons(authStore.auth.RoleId, current_message_status, authStore.auth.FriendRole)">
              <button @click="click_button(but.value, but.type_of_btn)" type="submit"
                class="inline-block rounded-lg bg-amber-400 hover:bg-amber-300 px-5 py-3 text-sm font-medium text-black w-full">
                {{ but.text }}
              </button>
            </div>
          </div>
          <van-action-sheet
          v-model:show="show"
          :actions="actions"
          cancel-text="Cancel"
          close-on-click-action
          @select="logout"
        />
        </div>
      </div>
    </div>
  </div>
  <van-floating-panel v-model:height="height" :anchors="anchors">
    <div class="px-4">
      <input
        type="number"
        v-model="friend_number"
        class="w-full rounded-lg border-gray-200 p-4 pe-12 text-sm shadow-sm border mt-2 mb-4"
        placeholder="Enter friend number"
      />
      <button
      @click="change_friend(friend_number, false, false)"
        type="submit"
        class="inline-block rounded-lg bg-amber-400 hover:bg-amber-300 px-5 py-3 text-sm font-medium text-black w-full"
      >
        Change Friend
      </button>
    </div>
  </van-floating-panel>
</template>
