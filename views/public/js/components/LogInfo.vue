<template>
  <div>
    <div
      class="overflow-x-hidden overflow-y-auto fixed inset-0 z-50 outline-none focus:outline-none justify-center items-center flex"
    >
      <div class="relative w-auto my-6 mx-auto max-w-6xl">
        <!--content-->
        <div
          class="border-0 rounded-lg shadow-lg relative flex flex-col w-full bg-white outline-none focus:outline-none"
        >
          <!--header-->
          <div
            class="flex items-start justify-between p-5 border-b border-solid border-gray-300 rounded-t"
          >
            <h4 class="text-2xl font-semibold">Request Details</h4>
            <button
              class="p-1 ml-auto bg-transparent border-0 text-black opacity-5 float-right text-3xl leading-none font-semibold outline-none focus:outline-none"
              @click.prevent="close"
            >
              <span
                class="bg-transparent text-black opacity-5 h-6 w-6 text-2xl block outline-none focus:outline-none"
              >
                ×
              </span>
            </button>
          </div>
          <!--body-->
          <div class="relative p-6 flex-auto">
            <p class="my-4 text-gray-600 text-lg leading-relaxed">
              <table class="min-w-full divide-y divide-gray-200">
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr><td class="px-6 py-4 whitespace-nowrap">Host :  {{info.host}}</td></tr> 
                  <tr><td class="px-6 py-4 whitespace-nowrap">Method :  {{info.method}}</td></tr> 
                  <tr><td class="px-6 py-4 whitespace-nowrap">Status :  {{info.code}}</td></tr> 
                  <tr><td class="px-6 py-4 whitespace-nowrap">Path :  {{info.path}}</td></tr> 
                  <tr><td class="px-6 py-4 whitespace-nowrap">Time Spent :  {{info.timespent}} ms</td></tr>   
                </tbody>
              </table>
              <div style='border-bottom: 2px solid #eaeaea'>
                <ul class='flex cursor-pointer'>
                    <li class='py-2 px-6 bg-white rounded-t-lg' v-bind:class="{'text-gray-500 bg-gray-200': openTab !== 1}"><a class="" v-on:click="toggleTabs(1)" >Payload</a></li>
                    <li class='py-2 px-6 bg-white rounded-t-lg' v-bind:class="{'text-gray-500 bg-gray-200': openTab !== 2}"><a class="" v-on:click="toggleTabs(2)">Headers</a></li>
                    <li class='py-2 px-6 bg-white rounded-t-lg' v-bind:class="{'text-gray-500 bg-gray-200': openTab !== 3}"> <a class="" v-on:click="toggleTabs(3)">Responses</a></li>
                </ul>
              </div>
              <div class="relative flex flex-col min-w-0 break-words bg-black w-full mb-6 shadow-lg rounded card-body">
                <div class="px-4 py-5 flex-auto">
                  <div class="tab-content tab-space text-white">
                    <div v-bind:class="{'hidden': openTab !== 1, 'block': openTab === 1}">
                      <p >
                        <pre>{{ info.requestbody | pretty }}</pre>
                      </p>
                    </div>
                    <div v-bind:class="{'hidden': openTab !== 2, 'block': openTab === 2}">
                      <p>
                         <pre>{{ info.headers | pretty }}</pre>
                      </p>
                    </div>
                    <div v-bind:class="{'hidden': openTab !== 3, 'block': openTab === 3}">
                      <p>
                         <pre>{{ info.responsebody | pretty }}</pre>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </p>
          </div>
          <!--footer-->
          <div
            class="flex items-center justify-end p-3 border-t border-solid border-gray-300 rounded-b"
          >
            <button
              class="text-red-500 bg-transparent border border-solid border-red-500 hover:bg-red-500 hover:text-white active:bg-red-600 font-bold uppercase text-sm px-6 py-3 rounded outline-none focus:outline-none mr-1 mb-1"
              type="button"
              style="transition: all 0.15s ease"
              @click.prevent="close"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
module.exports = {
  data() {
    return {
      openTab: 1
    }
  },
  props: {
    showing: {
      required: true,
      type: Boolean
    },
    info: {
      required: true,
      type: Object
    }
  },
  filters: {
    pretty: function(value) {
      return JSON.stringify(JSON.parse(value), null, 2);
    }
  },
  methods: {
    close() {
      this.$emit('close');
    },
    toggleTabs: function(tabNumber){
      this.openTab = tabNumber
    }
  }
};
</script>
<style scoped>
.card-body{
  background-color: black;
  padding:30px;
  height:300px;
  overflow:auto;
}
</style>