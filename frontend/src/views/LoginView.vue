<script setup lang="ts">
import { onMounted, ref,  } from 'vue';
import axios from "axios";
import {AuthStore} from "@/stores/authStore"
import {FCMStore} from "@/stores/fcm"
import logo from "@/assets/khanakhaneho_logo.png"
import router from '@/router';

let authStore = AuthStore()
let fcmStore = FCMStore()

let phone = ref("")
let password = ref("")

async function login(){
    try {
        let fcm = await fcmStore.GetFCMToke();
        if(fcm == ""){
            return
        }
        console.log(phone.value)
        let response = await axios.post(`${import.meta.env.VITE_HOST}/login`, {
            phone_number: `${phone.value}`,
            password: `${password.value}`,
            fcm: fcm,
        })


        if(response.status != 200){
            throw Error("Something went wrong");
        }

        let redirect = await authStore.setAuth(null, response.data.data.token)
        if(redirect){
            router.push({ name: 'dashboard'})
        }
    } catch (error) {
    }
}

onMounted(async()=>{
    let toke = await authStore.IsAuth()
    if(toke){
        router.push({ name: 'dashboard'})
    }
})
</script>

<template>
    <div class="w-screen h-screen flex flex-col justify-center ">
        <div class="mx-auto max-w-screen-xl px-4 py-16 sm:px-6 lg:px-8">
        <div class="mx-auto max-w-lg text-center">
          <img :src="logo" alt="" srcset="" class="w-full px-2">
          <h1 class="text-2xl font-bold sm:text-3xl mt-5">Login To Khanakhane ho?</h1>
      
          <p class="mt-4 text-gray-500">
            Log in and ask fellow users, "Khanakhane Ho?" to connect over meals.
          </p>
        </div>
      
        <div class="mx-auto mb-0 mt-8 max-w-md space-y-4">
          <div>
            <label for="email" class="sr-only">Number</label>
      
            <div class="relative">
              <input
                type="number"
                v-model="phone"
                class="w-full rounded-lg border-gray-200 p-4 pe-12 text-sm shadow-sm border"
                placeholder="Enter number"
              />
      
              <span class="absolute inset-y-0 end-0 grid place-content-center px-4">
                <span class="text-gray-400 font-semibold">+977</span>
              </span>
            </div>
          </div>
      
          <div>
            <label for="password" class="sr-only">Password</label>
      
            <div class="relative">
              <input
                v-model="password"
                type="password"
                class="w-full rounded-lg border border-gray-200 p-4 pe-12 text-sm shadow-sm"
                placeholder="Enter password"
              />
      
              <span class="absolute inset-y-0 end-0 grid place-content-center px-4">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-4 w-4 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                  />
                </svg>
              </span>
            </div>
          </div>
      
          <div class="flex items-center justify-between">
            <button
              @click="login()"
              type="submit"
              class="inline-block rounded-lg bg-amber-400 hover:bg-amber-300 px-5 py-3 text-sm font-medium text-black w-full"
            >
              Login / Register
            </button>
          </div>
        </div>
      </div>
    </div>
</template>