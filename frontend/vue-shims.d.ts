// declare module '*.vue' {
//   import { DefineComponent } from 'vue';
//   const component: DefineComponent<unknown, unknown, unknown>;
//   export default component;
// }

declare module '*.vue' {
  import type { DefineComponent } from 'vue';
  // TODO:
  // eslint-disable-next-line
  const component: DefineComponent<{}, {}, any>;
  export default component;
}
