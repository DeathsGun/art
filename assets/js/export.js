let form = document.getElementById("export");

getStartDate(document.getElementById("provider").value).then(date => {
    document.getElementById("date").value = date.toISOString().split("T")[0];
}).catch(e => console.error(e));

form.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    let body = {
        provider: document.getElementById("provider").value,
    };
    if (body.provider !== "PROVIDER_IHK") {
        body.date = new Date(document.getElementById("date").value).toISOString();
    }
    let resp = await fetch("/export", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "same-origin",
        body: JSON.stringify(body)
    });
    if (!resp.ok) {
        let message = document.getElementById("errorMessage");
        message.innerText = await resp.text();
        let errorModal = new bootstrap.Modal("#error");
        errorModal.show(undefined);
        return;
    }
    if (resp.status === 204) {
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
    let body = {
        provider: "PROVIDER_TEXT",
        date: new Date(document.getElementById("date").value).toISOString(),
    };
    let resp = await fetch("/export", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "same-origin",
        body: JSON.stringify(body)
    });
    if (!resp.ok) {
        let message = document.getElementById("errorMessage");
        message.innerText = await resp.text();
        let errorModal = new bootstrap.Modal("#error");
        errorModal.show(undefined);
        return;
    }
    let content = document.getElementById("previewContent");
    content.parentElement.classList.remove("d-none");
    content.value = await resp.text();
});

let provider = document.getElementById("provider");
provider.addEventListener("change", async () => {
    let prov = document.getElementById("provider").value;
    let date = document.getElementById("date");
    if (prov === "PROVIDER_IHK") {
        date.setAttribute("disabled", "true");
    } else {
        date.removeAttribute("disabled");
    }
    let startDate = await getStartDate(prov);
    date.value = startDate.toISOString().split("T")[0];
});

async function getStartDate(prov) {
    let resp = await fetch(`/export/start-date/${prov}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "same-origin",
    });
    if (!resp.ok) {
        console.error(await resp.text());
        return;
    }
    let body = await resp.json();
    return new Date(body.date);
}