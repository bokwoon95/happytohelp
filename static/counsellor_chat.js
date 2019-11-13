"use strict";

let conn;
let messages = [];
let roomId;
let displayname;

function appendLog(item) {
  let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
  messages.push(item);
  if (doScroll) {
    log.scrollTop = log.scrollHeight - log.clientHeight;
  }
}

function sendmsgOnSubmitOrEnter(event) {
  try {
    if (!conn) {
      return false;
    }
    let chatbox = document.querySelector("#chatbox");
    let chatlog = document.querySelector("#chatlog");
    let enter = event instanceof KeyboardEvent && event.key === "Enter" && !event.shiftKey;
    let mouse = event instanceof MouseEvent;
    if (mouse || enter) {
      if (!chatbox.value) {
        return false;
      }
      let doScroll = chatlog.scrollTop > chatlog.scrollHeight - chatlog.clientHeight - 1;
      conn.send(
        JSON.stringify({
          sender: displayname,
          message: chatbox.value,
        }),
      );
      if (doScroll) {
        chatlog.scrollTop = chatlog.scrollHeight - chatlog.clientHeight;
      }
      chat.scrollTop = 999;
      chatbox.value = "";
      return false;
    }
  } catch (err) {
    console.log(err);
  }
}

m.mount(document.querySelector("#chat"), {
  oninit: () => {
    try {
      if (!window["WebSocket"]) {
        messages.push("Your browser does not support WebSockets");
        return;
      }
      conn = new WebSocket(`ws://${document.location.host}/counsellor/chat/ws${document.location.search}`);
      conn.onclose = function(evt) {
        messages.push("Connection closed");
      };
      conn.onmessage = function(evt) {
        let responses = evt.data.split("\n");
        for (let response of responses) {
          let event = JSON.parse(response);
          let text =
            event.sender !== displayname
              ? `<b class="text-pink-700">Anon:</b> `
              : `<b class="text-teal-700">${event.sender}:</b> `;
          text += event.message;
          messages.push(text);
        }
        m.redraw();
      };
      let url = new URL(document.location.origin + document.location.pathname + document.location.search);
      roomId = url.searchParams.get("room");
      displayname = document.querySelector("#displayname") && document.querySelector("#displayname").innerHTML;
    } catch (err) {
      console.log(err);
    }
  },
  view: () =>
    m(
      "div.hth-bg-blue-300.min-h-screen",
      m(
        "div.hth-text-blue-200.pt-8.px-16.text-center.glacialindifference-reg.italic",
        "You can start chatting now. Thank you so much for helping out! This chat will be completely anonymous and no data will be saved.",
      ),
      m(
        "div.hth-bg-blue-300.overflow-auto.min-h-200px.p-8",
        { id: "chatlog" },
        messages.map(msg => m("div.mx-4.bg-gray-200.p-4.mb-8.rounded", m.trust(msg))),
      ),
      m(
        "form.px-12",
        m("textarea.w-full.bg-white.border-box.p-2.rounded", { id: "chatbox", onkeypress: sendmsgOnSubmitOrEnter }),
        m("button.p-2.bg-green-200.rounded", { onclick: sendmsgOnSubmitOrEnter }, "Submit"),
      ),
    ),
});
