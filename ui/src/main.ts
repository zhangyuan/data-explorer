import { createApp } from "vue";
import { createWebHashHistory, createRouter } from "vue-router";

import "./style.css";

import 'highlight.js/styles/stackoverflow-light.css'
import hljs from 'highlight.js/lib/core';
import javascript from 'highlight.js/lib/languages/sql';
import hljsVuePlugin from "@highlightjs/vue-plugin";
hljs.registerLanguage('sql', javascript);


import App from "@/App.vue";

import HomeView from "./components/HomeView.vue";
import NewIssue from "./components/issues/NewIssueView.vue";
import IssueListView from "./components/issues/IssueListView.vue";
import IssueView from "./components/issues/IssueView.vue";

const routes = [
  { path: "/", component: HomeView },
  { path: "/issues/new", component: NewIssue },
  { path: "/issues", component: IssueListView },
  { path: "/issues/:id", component: IssueView, name: "issue" },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

const app = createApp(App);
app.use(router);
app.use(hljsVuePlugin);
app.mount("#app");

hljs.highlightAll()
