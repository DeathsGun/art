let modals = document.querySelectorAll(".modal")

modals.forEach(modal => {
    let resetButton = document.querySelector(".btn-reset");
    resetButton.addEventListener("click", async event => {
        let button = event.target;
        let providerId = button.getAttribute("data-bs-provider");
        await resetConfig(providerId);
    });
    modal.addEventListener('show.bs.modal', async event => {
        let button = event.relatedTarget;
        let providerId = button.getAttribute("data-bs-provider-id");
        await loadConfig(providerId);
    });
});
async function loadConfig(name) {
    let resp = await fetch(`/config/${name}`, {
        credentials: "same-origin",
    });
    let modal = document.getElementById(name);
    let conf;
    if (!resp.ok) {
        if (resp.status !== 404) {
            let body = modal.querySelector(".modal-body");
            let alert = document.createElement("div");
            alert.classList.add("alert", "alert-danger");
            alert.setAttribute("role", "alert");
            alert.innerText = await resp.text();
            body.appendChild(alert);
            return;
        }
        conf = {
            server: "",
            username: "",
            password: "",
        }
    } else {
        conf = await resp.json();
    }

    let server = document.querySelector(`#${name}server`);
    if (server) {
        server.value = conf["server"];
    }
    // noinspection SpellCheckingInspection
    let username = document.querySelector(`#${name}username`);
    if (username) {
        username.value = conf["username"];
    }
    // noinspection SpellCheckingInspection
    let password = document.querySelector(`#${name}password`);
    if (password) {
        password.value = conf["password"];
    }
    console.log(conf);
    let department = document.querySelector(`#${name}department`);
    if (department) {
        department.value = conf["department"];
    }
    let instructor = document.querySelector(`#${name}instructor`);
    if (instructor) {
        instructor.value = conf["instructor"];
    }
    let form = document.getElementById(`${name}-form`);
    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        event.stopPropagation();
        if (!form.checkValidity()) {
            return;
        }
        let conf = {}
        if (server) {
            conf["server"] = server.value;
        }
        if (username) {
            conf["username"] = username.value;
        }
        if (password) {
            conf["password"] = password.value;
        }
        if (department) {
            conf["department"] = department.value;
        }
        if (instructor) {
            conf["instructor"] = instructor.value;
        }
        await saveConfig(name, conf);
    });
}

async function saveConfig(name, conf) {
    let resp = await fetch(`/config/${name}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "same-origin",
        body: JSON.stringify(conf),
    })
    let modal = document.getElementById(name);
    if (!resp.ok) {
        let body = modal.querySelector(".modal-body");
        for (let alert of modal.querySelectorAll("div.alert-danger")) {
            body.removeChild(alert);
        }
        let alert = document.createElement("div");
        alert.classList.add("alert", "alert-danger", "m-2");
        alert.setAttribute("role", "alert");
        alert.innerText = await resp.text();
        body.appendChild(alert);
        return;
    }
    const toast = new bootstrap.Toast(document.getElementById("saveToast"));
    toast.show();
    modal.querySelector(".btn-close").click();
}

async function resetConfig(providerId) {
    let resp = await fetch(`/config/${providerId}`, {
        method: "DELETE",
        credentials: "same-origin",
    })
    let modal = document.getElementById(providerId);
    if (!resp.ok) {
        let body = modal.querySelector(".modal-body");
        for (let alert of modal.querySelectorAll("div.alert-danger")) {
            body.removeChild(alert);
        }
        let alert = document.createElement("div");
        alert.classList.add("alert", "alert-danger", "m-2");
        alert.setAttribute("role", "alert");
        alert.innerText = await resp.text();
        body.appendChild(alert);
        return;
    }
    window.location.reload();
}
