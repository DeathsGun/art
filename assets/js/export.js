let form = document.getElementById("export");

document.getElementById("date").value = new Date().toISOString().split("T")[0];

form.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    let resp = await fetch("/export", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "same-origin",
        body: JSON.stringify({
            provider: document.getElementById("provider").value,
            date: new Date(document.getElementById("date").value).toISOString(),
        })
    });
    if (!resp.ok) {
        console.error(await resp.text());
        return;
    }
    let url = window.URL.createObjectURL(await resp.blob());
    let download = document.createElement("a");
    download.href = url;
    download.target = "_blank";
    download.download = resp.headers.get("File-Name");
    download.click();
});

let previewButton = document.getElementById("preview");
previewButton.addEventListener("click", async () => {
    let resp = await fetch("/export", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "same-origin",
        body: JSON.stringify({
            provider: "PROVIDER_TEXT",
            date: new Date(document.getElementById("date").value).toISOString(),
        })
    });
    let content = document.getElementById("previewContent");
    content.parentElement.classList.remove("d-none");
    content.value = await resp.text();
});