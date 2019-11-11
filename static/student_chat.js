"use strict";

let conn;
const messages = ["hello world"];

function sendmsg() {
  try {
    if (!conn) {
      return false;
    }
    const chatbox = document.querySelector("#chatbox");
    if (!chatbox.value) {
      return false;
    }
    conn.send(chatbox.value);
    chatbox.value = "";
    return false;
  } catch (err) {
    console.log(err);
  }
}

m.mount(document.querySelector("#chat"), {
  oninit: () => {
    if (window["WebSocket"]) {
      conn = new WebSocket(`ws://${document.location.host}/student-chat/ws`);
      conn.onclose = function(evt) {
        messages.push("Connection closed");
      };
      conn.onmessage = function(evt) {
        var responses = evt.data.split("\n");
        console.log(responses);
        for (const response of responses) {
          messages.push(response);
        }
        m.redraw();
      };
    } else {
      messages.push("Your browser does not support WebSockets");
    }
  },
  view: () =>
    m(
      "div.bg-blue-800.min-h-screen",
      m("div.bg-gray-100.overflow-auto.min-h-200px.p-8", messages.map(msg => m("p", msg))),
      m(
        "form",
        m("textarea.w-full.bg-gray-200", { id: "chatbox" }),
        m("button.p-2.bg-green-200", { onclick: sendmsg }, "Submit"),
      ),
    ),
});
