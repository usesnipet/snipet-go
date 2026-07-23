export type MockSession = {
	id: string;
	metadata: { name?: string; [key: string]: unknown };
};

function createInitialSessions(): MockSession[] {
	return [
		{ id: crypto.randomUUID(), metadata: { name: "Welcome chat" } },
		{ id: crypto.randomUUID(), metadata: { name: "Agent playground" } },
		{ id: crypto.randomUUID(), metadata: { name: "Knowledge test" } },
	];
}

class MockSessionsStore {
	sessions = $state<MockSession[]>(createInitialSessions());

	create(name?: string): MockSession {
		const session: MockSession = {
			id: crypto.randomUUID(),
			metadata: { name: name?.trim() || "New chat" },
		};
		this.sessions = [session, ...this.sessions];
		return session;
	}

	rename(id: string, name: string): void {
		const trimmed = name.trim() || "New chat";
		this.sessions = this.sessions.map((session) =>
			session.id === id
				? { ...session, metadata: { ...session.metadata, name: trimmed } }
				: session,
		);
	}

	delete(id: string): void {
		this.sessions = this.sessions.filter((session) => session.id !== id);
	}

	filtered(query: string): MockSession[] {
		const q = query.trim().toLowerCase();
		if (!q) return this.sessions;
		return this.sessions.filter((session) => {
			const name = String(session.metadata.name ?? "").toLowerCase();
			return name.includes(q);
		});
	}
}

export const mockSessions = new MockSessionsStore();
