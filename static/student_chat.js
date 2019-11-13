"use strict";

let conn; // variable that holds the websocket connection with the server
let messages = []; // list of messages to display
let roomId;

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
          sender: roomId,
          message: chatbox.value,
        }),
      );
      if (doScroll) {
        chatlog.scrollTop = chatlog.scrollHeight - chatlog.clientHeight;
      }
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
      conn = new WebSocket(`ws://${document.location.host}/student/chat/ws${document.location.search}`);
      conn.onclose = function(evt) {
        messages.push("Connection closed");
      };
      conn.onmessage = function(evt) {
        let responses = evt.data.split("\n");
        for (let response of responses) {
          let event = JSON.parse(response);
          let text =
            event.sender === roomId
              ? `<b class="text-teal-700">You:</b> `
              : `<b class="text-pink-700">${event.sender}:</b> `;
          text += event.message;
          messages.push(text);
        }
        m.redraw();
      };
      let url = new URL(document.location.origin + document.location.pathname + document.location.search);
      roomId = url.searchParams.get("room");
    } catch (err) {
      console.log(err);
    }
  },
  view: () => {
    return messages.length == 0
      ? m(
          "div",
          m(
            "div.text-center.text-5xl.hth-text-blue-300.px-16",
            "Thank you! Please give our volunteers some time to connect with you",
          ),
          m("img.m-auto", { src: "/static/img/HappyToHelp_Logo.svg", height: 500, width: 500 }),
        )
      : m(
          "div.hth-bg-blue-300.min-h-screen",
          m(
            "div.hth-text-blue-200.pt-8.px-16.text-center.glacialindifference-reg.italic",
            "You can start chatting now. This chat will be completely anonymous and no data will be saved. Please be assured that all our volunteers are your fellow students.",
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
        );
  },
});
