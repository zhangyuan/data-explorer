import { createApp } from "vue";
import { createWebHashHistory, createRouter } from "vue-router";

import "./style.css";
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
app.mount("#app");
