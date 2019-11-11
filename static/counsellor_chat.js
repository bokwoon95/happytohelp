"use strict";

m.mount(document.querySelector("#log"), {
  view: () => m("p", "hello world!"),
})
