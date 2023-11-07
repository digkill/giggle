import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';
import VueNativeSock from 'vue-native-websocket';

const BACKEND_URL = 'http://localhost:8080';
const PUSHER_URL = 'ws://localhost:8080/pusher';

const SET_GIGGLES = 'SET_GIGGLES';
const CREATE_GIGGLE = 'CREATE_GIGGLE';
const SEARCH_SUCCESS = 'SEARCH_SUCCESS';
const SEARCH_ERROR = 'SEARCH_ERROR';

const MESSAGE_GIGGLE_CREATED = 1;

Vue.use(Vuex);

const store = new Vuex.Store({
    state: {
        giggles: [],
        searchResults: [],
    },
    mutations: {
        SOCKET_ONOPEN(state, event) {
        },
        SOCKET_ONCLOSE(state, event) {
        },
        SOCKET_ONERROR(state, event) {
            console.error(event);
        },
        SOCKET_ONMESSAGE(state, message) {
            switch (message.kind) {
                case MESSAGE_GIGGLE_CREATED:
                    this.commit(CREATE_GIGGLE, { id: message.id, body: message.body });
            }
        },
        [SET_GIGGLES](state, giggles) {
            state.giggles = giggles;
        },
        [CREATE_GIGGLE](state, giggle) {
            state.giggles = [giggle, ...state.giggles];
        },
        [SEARCH_SUCCESS](state, giggles) {
            state.searchResults = giggles;
        },
        [SEARCH_ERROR](state) {
            state.searchResults = [];
        },
    },
    actions: {
        getGiggles({ commit }) {
            axios
                .get(`${BACKEND_URL}/giggles`)
                .then(({ data }) => {
                    commit(SET_GIGGLES, data);
                })
                .catch((err) => console.error(err));
        },
        async createGiggle({ commit }, giggle) {
            const { data } = await axios.post(`${BACKEND_URL}/giggles`, null, {
                params: {
                    body: giggle.body,
                },
            });
        },
        async searchGiggles({ commit }, query) {
            if (query.length === 0) {
                commit(SEARCH_SUCCESS, []);
                return;
            }
            axios
                .get(`${BACKEND_URL}/search`, {
                    params: { query },
                })
                .then(({ data }) => commit(SEARCH_SUCCESS, data))
                .catch((err) => {
                    console.error(err);
                    commit(SEARCH_ERROR);
                });
        },
    },
});

Vue.use(VueNativeSock, PUSHER_URL, { store, format: 'json' });

store.dispatch('getGiggles');

export default store;