if (localStorage.getItem("school")) {
    document.getElementById("school").style.display = "none";
    document.getElementById("account").style.display = "block";
    document.getElementById("schoolField").value = localStorage.getItem("school");
} else {
    document.getElementById("school").style.display = "block";
    document.getElementById("account").style.display = "none";
}

async function search() {
    let query = document.getElementById("schoolSearch").value;
    if (query.length < 3) {
        return;
    }
    let response = await fetch("/untis/search", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            search: query
        }),
    });
    if (!response.ok) {
        console.log(response.text());
        return;
    }
    let schools = await response.json();
    if (document.getElementById("schools")) {
        document.getElementById("schools").remove();
    }
    let list = document.createElement("div");
    list.id = "schools";
    list.classList.add("list-group", "text-start");
    for (const school of schools) {
        let item = document.createElement("a");
        item.classList.add("list-group-item", "list-group-item-action")
        item.id = school.loginName;
        item.href = "#";
        item.onclick = () => {
            localStorage.setItem("school", school.loginName);
            document.getElementById("selectSchool").removeAttribute("disabled");
            let children = document.getElementById("schools").children;
            for (const child of children) {
                child.classList.remove("active");
            }
            document.getElementById(school.loginName).classList.add("active");
        }

        let container = document.createElement("div");
        container.classList.add("d-flex", "w-100", "justify-content-between");

        let title = document.createElement("h5");
        title.classList.add("mb-1");
        title.innerText = school.displayName;

        let address = document.createElement("p");
        address.classList.add("mb-1");
        address.innerText = school.address;

        container.appendChild(title);
        item.appendChild(container);

        item.appendChild(address);
        list.appendChild(item);
    }
    document.getElementById("school").insertBefore(list, document.getElementById("selectSchool"));
}

async function setSchool() {
    if (!localStorage.getItem("school")) {
        return;
    }
    document.getElementById("school").style.display = "none";
    document.getElementById("account").style.display = "block";
}

function resetSchool() {
    localStorage.removeItem("school");
    window.location.reload();
}
