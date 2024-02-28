function makeFfs(selector, requestKey = "value") {
	function doNode(form, n) {
		async function upload() {
			const body = {};
			body[requestKey] = form.dataset.value === "true";
			body["solved"] = form.dataset.solved === "true";
			const resp = await fetch(form.action, {
				method: form.method || "POST",
				body: JSON.stringify(body),
				headers: { "Content-Type": "application/json" },
			}).catch(() => null);
			const j = await resp?.json().catch((err) => {
				console.warn(err);
				return false;
			});
			form.dataset.confirmed = !!j;
			console.log({ j });
			return j;
		}

		function upd(value, solved = true) {
			form.dataset.value = value;
			form.dataset.solved = solved;
			upload();
		}

		form.onclick = () => !!upd(true);
		form.oncontextmenu = () => !!upd(false);
		form.ondblclick = () => !!upd(false, false);
	}

	Array.from(document.querySelectorAll(selector)).forEach(doNode);
}
