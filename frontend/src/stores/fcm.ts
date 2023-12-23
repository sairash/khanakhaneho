import { ref, type Ref } from 'vue'
import { defineStore } from 'pinia'
import type { FirebaseApp } from 'firebase/app'
import { Capacitor } from '@capacitor/core';
import { initializeApp } from "firebase/app";
import { getMessaging, getToken } from "firebase/messaging";
import {LocalNotifications} from '@capacitor/local-notifications';

import {
    PushNotifications, type ActionPerformed, type Token, type PushNotificationSchema,
  } from '@capacitor/push-notifications';


export const FCMStore = defineStore('fcm', () => {
    const app = ref<FirebaseApp>()
    const fcm = ref("")
  
    async function GetNewFCMToken(): Promise<string> {
        try {
            let app_as_fireapp = app.value as FirebaseApp
            const messaging = await getMessaging(app_as_fireapp);
            let curren_token = await getToken(messaging, {vapidKey: "BBZX-LdmArn6URNjBz3Sgboj2pbEz8lkF5xkSqAIT3MJ4muJ5l1BQWujiUCkzJ_3Wi8yQO0h_J3bywZZb9hXxUI"})
            if(curren_token){
                return curren_token
            }
            throw Error("Notification Perminssion not allowed.")
        } catch (error) {
            console.log(error)
            return ""
        }
    }
  
    async function GetFCMToke(): Promise<string> {
        if(fcm.value == ""){
            let new_token = await GetNewFCMToken()
            if(new_token != ""){
                fcm.value = new_token
            }
        }
        return fcm.value
    }

    async function SetApp(){
        if (Capacitor.getPlatform() === 'web') {
            const firebaseConfig = {
                apiKey: "AIzaSyDaCMHY6exJKeFjX88XBlLPolLZMRI-Fic",
                authDomain: "khanakhaneho-sai.firebaseapp.com",
                projectId: "khanakhaneho-sai",
                storageBucket: "khanakhaneho-sai.appspot.com",
                messagingSenderId: "499598278288",
                appId: "1:499598278288:web:0465078fb8de22ede76215"
            };
              
            app.value = initializeApp(firebaseConfig);
            await GetFCMToke()
        }else{
            PushNotifications.requestPermissions().then(result => {
                if (result.receive === 'granted') {
                    PushNotifications.register();
                } 
            });

            PushNotifications.addListener('registration', (token: Token) => {
                fcm.value = token.value;
            });


            PushNotifications.addListener(
                'pushNotificationReceived',
                (notification: PushNotificationSchema) => {
                    LocalNotifications.schedule(({
                        notifications:[{
                            id: 111,
                            title: notification.title as string,
                            body: notification.body as string,

                        }] 
                    }));
                    location.reload(); 
                },
              );

            // PushNotifications.addListener(
            // 'pushNotificationActionPerformed',
            // (notification: ActionPerformed) => {
            //     alert('Push action performed: ' + JSON.stringify(notification));
            // },
            // );
        }

    }
  
    return { app, fcm, GetFCMToke, SetApp }
  })