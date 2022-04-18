window.recorder = {
	events: [],
	rrweb: undefined,
	runner: undefined,
	session: {
		genId(length) {
			const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
			let result = "";
			const charactersLength = characters.length;
			for (let i = 0; i < length; i++) {
				result += characters.charAt(Math.floor(Math.random() * charactersLength));
			}
			return result;
		},
		get() {
			let session = window.sessionStorage.getItem('rrweb');
			if (session) return JSON.parse(session);
			session = {
				id: window.recorder.session.genId(64),
				user: { id: window.recorder.session.genId(64) },
				clientId: 'default'
			};
			window.sessionStorage.setItem('rrweb', JSON.stringify(session));
			return session;
		},
		save(data) {
			const session = window.recorder.session.get();
			window.sessionStorage.setItem('rrweb', JSON.stringify(Object.assign({}, session, data)));
		},
		clear() {
			window.sessionStorage.removeItem('rrweb')
		}
	},
	setUser: function({ id, email, name }) {
		const session = window.recorder.session.get();
		session.user = { id, email, name };
		window.recorder.session.save(session)
		return window.recorder;
	},
	setClientId(id) {
		const session = window.recorder.session.get();
		session.clientId = id;
		window.recorder.session.save(session)
		return window.recorder;
	},
	stop() {
		clearInterval(window.recorder.runner);
	},
	start() {
		window.recorder.runner = setInterval(function save() {
			const session = window.recorder.session.get();
			fetch('{{ .URL }}/sessions', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(Object.assign({}, { events: window.recorder.events }, session)),
			});
			window.recorder.events = []; // cleans-up events for next cycle
		}, 5 * 1000);
	},
	close() {
		clearInterval();
		window.recorder.session.clear();
	}
};
new Promise((resolve, reject) => {
	const script = document.createElement('script');
	script.src = 'https://cdn.jsdelivr.net/npm/rrweb@latest/dist/rrweb.min.js';
	script.addEventListener('load', resolve);
	script.addEventListener('error', e => reject(e.error));
	document.head.appendChild(script);
}).then(() => {
	window.recorder.rrweb = rrweb;
	// TODO: This should be optimised
	rrweb.record({
		emit(event) {
			window.recorder.events.push(event);
		}
	});
	window.recorder.start();
}).catch(console.err);