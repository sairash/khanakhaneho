import { ref } from 'vue'
import { defineStore } from 'pinia'
import { Preferences } from '@capacitor/preferences'
import axios from 'axios'

interface user_data {
    Id: number,
    PhoneNumber: number,
    RoleId: number,
    FriendPhone: number,
    FriendRole: number
}


export const AuthStore = defineStore('auth', () => {
  const auth = ref<user_data>({
    Id: 0,
    PhoneNumber: 0,
    RoleId: 0,
    FriendPhone: 0,
    FriendRole: 0
  })
  const token = ref("")

  const IsAuth = async()=>{
    // return true
    let x = auth.value as user_data;
    return await checkAuth(x)
  }

  const checkAuth = async (value: user_data) => {
    if(value.Id == 0){
        if(token.value == ""){
            const { value } = await Preferences.get({ key: 'token' });
            if(value == null || value == ""){
                return false;
            }
            token.value = value;
        }
        return check_me(token.value);
    }
    return true;
  };


  const check_me = async (token: string) => {
    try {
        let response = await axios.post(`${import.meta.env.VITE_HOST}/me`,{}, {
          headers:{
            'Authorization': `Bearer ${token}`
          }
        })
        if(response.status != 200){
            throw new Error(`Login Again!`);
        }
        return await setAuth(response.data.data.user as user_data)
    } catch (error) {
        await Preferences.remove({key: "token"})
        return false
    }
  };
  

  async function setAuth(auth_from_resp: user_data|null, token_from_resp: string|null=null) : Promise<boolean>{
    if(token_from_resp != null){
        await Preferences.set({
            key: 'token',
            value: token_from_resp,
        });
        token.value = token_from_resp;
        return IsAuth()
    }
    auth.value = auth_from_resp as user_data;
    return true
  }


  async function deleteAuth():Promise<boolean> {
    auth.value = {
      Id: 0,
      PhoneNumber: 0,
      RoleId: 0,
      FriendPhone: 0,
      FriendRole: 0
    }
    await Preferences.remove({
      key: 'token',
    });
    location.reload(); 
    return true
  }

  return { auth, token, IsAuth, setAuth, deleteAuth }
})