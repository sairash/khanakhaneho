import { ref, type Ref } from 'vue'
import { defineStore } from 'pinia'
import sad from "@/assets/1F620.svg"
import thinking from "@/assets/1F914.svg"
import happy from "@/assets/1F917.svg"
import tongue_bite from "@/assets/tongue_bite.svg"


interface Buttons {
    text: string;
    value: string;
    type_of_btn: string;
}

interface Question {
    starter: string;
    ender: string;
    image: string;
    extra: string;
    buttons: Buttons[]
  }
  
  // Define the questions object using the interface
  interface Questions {
    [key: string]: {
      [key: string]: {
        [subKey: string]: Question;
      };
    };
  }

export const DashboardStore = defineStore('dashboard', () => {
    const questions = ref<Questions>({
        son:{
          friend:{
            not_asked:{
              starter: "Add a ",
              ender: " first to ask or be asked",
              image: tongue_bite,
              buttons: [],
              extra: ""
            }
          },
          mom:{
            not_asked: {
              starter: "Wait for your ",
              image: thinking,
              extra: "",
              ender: " to ask question.",
              buttons: []
            },
            question_will_you_eat_today: {
              starter: "Answer Your ",
              image: thinking,
              extra: "",
              ender: " before 3 hrs.",
              buttons: [{
                  text: "Yes (Ahh)!",
                  value: "yes_i_will_eat_today",
                  type_of_btn: "answer"
              },{
                  text: "No (Aha)!",
                  value: "no_i_will_not_eat_today",
                  type_of_btn: "answer"
              }]
            },
            answer_yes_i_will_eat_today: {
              starter: "You told ",
              image: happy,
              extra: "confetti",
              ender: " that you will have the next meal!",
              buttons: []
            },
            answer_no_i_will_not_eat_today: {
              starter: "You told ",
              extra: "sad",
              image: sad,
              ender: " that you will not have the next meal!",
              buttons: []
            },
          }
        },
        mom:{
          friend:{
            not_asked:{
              starter: "Add a ",
              ender: " first to ask or be asked",
              image: tongue_bite,
              buttons: [],
              extra: ""
            }
          },
          son:{
            not_asked: {
              starter: "Ask your",
              image: thinking,
              extra: "",
              ender: " if they want to eat today?",
              buttons: [{
                  text: "Will you eat today?",
                  value: "will_you_eat_today",
                  type_of_btn: "question"
              }]
            },
            question_will_you_eat_today: {
              starter: "Wait for your ",
              extra: "",
              image: tongue_bite,
              ender: " to answer!",
              buttons: []
            },
            answer_yes_i_will_eat_today: {
              starter: "Your ",
              image: happy,
              extra: "confetti",
              ender: " told that you will have the next meal!",
               buttons: [{
                  text: "Ask again will you eat today?",
                  value: "will_you_eat_today",
                  type_of_btn: "question"
              }]
            },
            answer_no_i_will_not_eat_today: {
              starter: "Your ",
              image: sad,
              extra: "sad",
              ender: " told that you will not have the next meal!",
               buttons: [{
                  text: "Ask again will you eat today?",
                  value: "will_you_eat_today",
                  type_of_btn: "question"
              }]
            },
          }
        }
      })
    const roles = ref(["friend", "", "son", "mom"])
  
  
    function get_question(role_id: number, needed_value: string, current_message_status: string, friend_role: number): string {
        switch (needed_value) {
            case "starter":
                return questions.value[roles.value[role_id]][roles.value[friend_role]][current_message_status].starter
            case "ender":
                return questions.value[roles.value[role_id]][roles.value[friend_role]][current_message_status].ender
            case "role":
                return roles.value[role_id]
            case "image":
                return questions.value[roles.value[role_id]][roles.value[friend_role]][current_message_status].image
            case "extra":
              return questions.value[roles.value[role_id]][roles.value[friend_role]][current_message_status].extra
            case "bg":
              switch (questions.value[roles.value[role_id]][roles.value[friend_role]][current_message_status].extra) {
                case "confetti":
                  return "bg-green-300"
                case "sad":
                  return "bg-red-300"
                default:
                  return "bg-white"
              }
        }
        return ""
    }

    function get_buttons(role_id: number, current_message_status: string, friend_role: number): Buttons[]{
        return questions.value[roles.value[role_id]][roles.value[friend_role]][current_message_status].buttons
    }

   
  
    return { questions, roles, get_question, get_buttons }
  })