import { Routine, APIResponse } from "../types/routine";

export default async function Home() {
  const res = await fetch("http://backend:8080/api/v1/routines", {
    cache: "no-store",
  });
  const json: APIResponse<Routine[]> = await res.json();
  const routines = json.data;

  return (
    <main className="p-8">
      <h1 className="text-2xl font-bold mb-6">My Routines</h1>
      <ul className="space-y-4">
        {routines.map((r) => (
          <li key={r.id} className="p-4 border rounded-lg shadow-sm bg-white">
            <h2 className="font-semibold text-lg">{r.title}</h2>
            <p className="text-gray-600">Interval: {r.interval}</p>
          </li>
        ))}
      </ul>
    </main>
  );
}
