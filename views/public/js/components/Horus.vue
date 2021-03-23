<template>
  <div>
    <div class="-mx-4 sm:-mx-8 px-4 sm:px-8 py-4 overflow-x-auto">
      <div class="inline-block min-w-full shadow rounded-lg overflow-hidden">
        Live Monitoring :
        <button v-if="status === 'connected'"
                class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">
          Connected
        </button>
        <button v-if="status === 'disconnected'" aria-hidden
                class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">
          Disconnected
        </button>
        <table class="min-w-full leading-normal">
          <thead>
          <tr>
            <th
                class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider"
            >
              Verb
            </th>
            <th
                class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider"
            >
              Path
            </th>
            <th
                class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider"
            >
              Status
            </th>
            <th
                class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider"
            >
              Duration
            </th>
            <th
                class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider"
            >
              Created At
            </th>
            <th
                class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider"
            >
              👁️
            </th>
          </tr>
          </thead>
          <tbody v-for="log in logs" v-bind:key="log.id" v-bind:log="log">
          <!-- <log-view v-for="log in logs" v-bind:key="log.id" v-bind:log="log"/> -->
          <tr>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              <p class="text-gray-900 whitespace-no-wrap">{{ log.method }}</p>
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              <p class="text-gray-900 whitespace-no-wrap">{{ log.path }}</p>
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
                <span
                    class="relative inline-block px-3 py-1 font-semibold text-green-900 leading-tight"
                >
                  <span
                      aria-hidden
                      class="absolute inset-0 opacity-50 rounded-full"
                      v-bind:class="{'bg-red-200': log.code !== 200, 'bg-green-200' : log.code === 200}"
                  ></span>
                  <span class="relative">{{ log.code }}</span>
                </span>
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              {{ log.timespent }} ms
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              {{ log.CreatedAt }}
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              <a v-on:click="toggleModal(log.ID)"> 👁️ </a>
            </td>
          </tr>
          </tbody>
        </table>
        <br/>
        <br/>
        <br/>
      </div>
    </div>
    <div v-if="showLog">
      <log-info :showing="showLog" :info="selectedLog" @close="showLog = false"/>
    </div>
    <div v-if="showLog" class="opacity-25 fixed inset-0 z-40 bg-black"></div>
  </div>
</template>

<script>
module.exports = {
  data() {
    return {
      isLoading: false,
      showLog: false,
      selectedLog: {},
      logs: [],
      connection: null,
      status: 'disconnected'
    };
  },
  components: {
    "log-view": httpVueLoader("public/js/components/Logs.vue"),
    "log-info": httpVueLoader("public/js/components/LogInfo.vue"),
  },
  mounted() {
    this.initiate()
    this.scroll();
  },
  methods: {
    initiate() {
      fetch("./logs?lastID=0")
          .then((response) => response.json())
          .then((data) => (this.logs = data));

      this.connection = new WebSocket("ws://" + document.location.host + "/ws");

      this.connection.onclose = () => {
        this.status = "disconnected"
      }

      this.connection.onopen = () => {
        this.status = "connected"
      }

      this.connection.onmessage = (evt) => {
        const js = JSON.parse(evt.data);
        this.logs.unshift(js)
      }
    },
    toggleModal(id) {
      this.selectedLog = this.logs.find(obj => {
        return obj.ID === id
      });
      this.showLog = !this.showLog;
    },
    scroll() {
      window.onscroll = () => {
        let bottomOfWindow = document.documentElement.scrollTop + window.innerHeight === document.documentElement.offsetHeight;

        if (bottomOfWindow) {
          let lastItem = this.logs.slice(-1);

          let url = "./logs"

          if (lastItem[0].ID !== undefined) {
            url += "?lastID=" + lastItem[0].ID

            fetch(url)
                .then((response) => response.json())
                .then((data) => (
                    this.logs = this.logs.concat(data)
                ));
          }
        }
      };
    },
  },
};
</script>

<style scoped>
.modal {
  transition: opacity 0.25s ease;
}

body.modal-active {
  overflow-x: hidden;
  overflow-y: visible !important;
}
</style>