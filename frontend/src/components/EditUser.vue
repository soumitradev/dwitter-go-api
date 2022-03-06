<template>
  <div class="flex flex-col bg-neutral-99 max-w-xl pl-4 pb-4">
    <div class="flex flex-row justify-between">
      <div class="flex flex-row mt-4">
        <div class="w-32 h-32 rounded-full overflow-hidden">
          <div
            class="hover:bg-neutral-50 bg-opacity-0 hover:bg-opacity-30 rounded-full hover:block w-32 h-32 absolute group"
            @click="toggleShow"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-full w-10 fill-current mx-auto group-hover:opacity-100 opacity-0 stroke-neutral-99 fill-transparent"
              viewBox="0 0 24 24"
              fill="none"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"
              />
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>

            <!-- <input
              type="file"
              class="opacity-0 absolute left-0 top-0 h-32 w-32 rounded-full cursor-pointer"
              @change="toggleShow"
              accept="image/png, image/jpeg, image/gif"
            />-->
          </div>

          <img :src="imgURL ? imgURL : pfpURL" />
        </div>
        <div class="flex flex-col ml-4 text-xl self-center">
          <input
            class="bg-neutral-99 mt-2 text-left font-bold text-neutral-10 border-neutralVariant-50 border-opacity-50 border-x-0 border-t-0 border-b-2 resize-none appearance-none outline-none"
            name="text"
            oninput="this.style.height = '';this.style.height = (this.scrollHeight + 1) + 'px'"
            rows="1"
            v-model="nameVar"
            autofocus
          />
          <span class="text-left text-neutralVariant-50">@{{ username }}</span>
        </div>
      </div>
      <button
        class="bg-primary-40 rounded-full m-4 h-fit px-2 flex flex-row text-primary-100"
        @click="$emit('saveEditsToUser', consolidateData())"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 m-1 fill-current"
          viewBox="0 0 20 20"
        >
          <path
            fill-rule="evenodd"
            d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
            clip-rule="evenodd"
          />
        </svg>
      </button>
    </div>

    <textarea
      class="bg-neutral-99 mt-4 text-left px-2 text-neutral-10 mx-2 border-neutralVariant-50 border-opacity-50 border-x-0 border-t-0 border-b-2 resize-none appearance-none outline-none"
      name="text"
      oninput="this.style.height = '';this.style.height = (this.scrollHeight + 5) + 'px'"
      rows="1"
      v-model="bioVar"
    ></textarea>

    <div class="flex flex-row mt-2 text-neutralVariant-50">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5 mx-1 my-1 fill-current"
        viewBox="0 0 24 24"
      >
        <path d="M0 0h24v24H0V0z" fill="none" />
        <path
          d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm-.4 4.25l-7.07 4.42c-.32.2-.74.2-1.06 0L4.4 8.25c-.25-.16-.4-.43-.4-.72 0-.67.73-1.07 1.3-.72L12 11l6.7-4.19c.57-.35 1.3.05 1.3.72 0 .29-.15.56-.4.72z"
        />
      </svg>
      <input
        class="bg-neutral-99 text-left text-neutral-10 border-neutralVariant-50 border-opacity-50 border-x-0 border-t-0 border-b-2 resize-none appearance-none w-full mx-2 outline-none"
        name="text"
        oninput="this.style.height = '';this.style.height = (this.scrollHeight + 1) + 'px'"
        rows="1"
        v-model="emailVar"
      />
    </div>
  </div>

  <my-upload
    :field="'img' + viewUser"
    @crop-success="cropSuccess"
    v-model="show"
    :width="240"
    :height="240"
    img-format="png"
    langType="en"
    :noSquare="true"
    :noRotate="false"
  ></my-upload>
</template>

<script>
import { ref } from 'vue';
import myUpload from 'vue-image-crop-upload';

export default {
  name: "EditUser",
  methods: {
    cropSuccess: function (imgURL, _) {
      this.imgURL = imgURL;
    },
    consolidateData: function () {
      return {
        username: this.username,
        name: this.nameVar,
        email: this.emailVar,
        bio: this.bioVar,
        pfpURL: this.imgURL ? this.imgURL : this.pfpURL,
        pfpUpdated: this.imgURL ? true : false,
      }
    },
    // updateFile: function (event) {
    //   console.log(event.target.files[0]);
    //   let valid = this.validateFile(event.target.files[0]);
    //   if (valid) {
    //     this.pfpFile = event.target.files[0];
    //   } else {
    //     this.openModal();
    //   }
    // },
    // validateFile: function (file) {
    //   return (file.size <= (8 << 20));
    // },
    // getURL: function (file) {
    //   console.log(file);
    //   return URL.createObjectURL(file);
    // },
  },
  props: {
    username: {
      type: String,
    },
    name: {
      type: String,
    },
    email: {
      type: String,
    },
    bio: {
      type: String,
    },
    pfpURL: {
      type: String,
    },
    viewUser: {
      type: String,
    }
  },
  data() {
    return {
      nameVar: this.name,
      emailVar: this.email,
      bioVar: this.bio,
    }
  },
  components: {
    'my-upload': myUpload,
  },
  emits: ['saveEditsToUser'],
  setup() {
    const isOpen = ref(false);
    const show = ref(false);
    const imgURL = ref("");

    return {
      isOpen,
      show,
      imgURL,
      toggleShow() {
        this.show = !this.show;
      },
    }
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
