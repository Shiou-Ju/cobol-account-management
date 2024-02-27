// TODO: this file sometime cause tsconfig showing error:
// No inputs were found in config file '/Users/bamboo/Repos/cobol-account-management/frontend/tsconfig.json'. Specified 'include' paths were '["src/**/*.ts","src/**/*.vue","src/vue-shims.d.ts"]' and 'exclude' paths were '["node_modules"]'.ts

declare module '*.vue' {
  import { DefineComponent } from 'vue';
  // TODO: not sure using unknown is good practice
  const component: DefineComponent<unknown, unknown, unknown>;
  export default component;
}
