const { VITE_HIDE_HOME } = import.meta.env;
const Layout = () => import("@/layout/index.vue");

export default {
  path: "/",
  name: "Contest",
  component: Layout,
  redirect: "/contest",
  meta: {
    title: "竞赛",
    rank: 0
  },
  children: [
    {
      path: "/contest",
      name: "Contest",
      component: () => import("@/views/contest/index.vue"),
      meta: {
        title: "竞赛",
        showLink: VITE_HIDE_HOME === "true" ? false : true
      }
    }
  ]
} satisfies RouteConfigsTable;
