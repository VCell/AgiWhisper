const serverResponse = document.getElementById('serverResponse');

function sendAudio() {
    const fileInput = document.getElementById('audioFile');
    if (fileInput.files.length == 0) {
        alert('Please select an audio file first.');
        return;
    }
    const file = fileInput.files[0];
    const reader = new FileReader();

    reader.onload = function (e) {
        const audioData = new Uint8Array(e.target.result);
        const chunkSize = 8000; // 8k chunk size
        var loc = window.location, ws_url;
        if (loc.protocol === "https:") {
            ws_url = "wss:";
        } else {
            ws_url = "ws:";
        }
        ws_url += "//" + loc.host + "/talk_manual";
        const ws = new WebSocket(ws_url);
        ws.onopen = function () {
            for (let i = 0; i < audioData.length; i += chunkSize) {
                const chunk = audioData.slice(i, i + chunkSize);
                const base64Chunk = btoa(String.fromCharCode.apply(null, chunk));
                const frame = JSON.stringify({ audio: base64Chunk });
                ws.send(frame);
            }
            ws.send(JSON.stringify({ audio: "", action: "ask" }));
        };

        ws.onmessage = function (event) {
            serverResponse.textContent += `\n${event.data}`;
        };

        ws.onerror = function (event) {
            alert('WebSocket error: ' + event.message);
        };

        ws.onclose = function (event) {
            console.log('WebSocket connection closed');
        };
    };

    reader.onerror = function () {
        alert('Failed to read the file.');
    };

    reader.readAsArrayBuffer(file);
}

