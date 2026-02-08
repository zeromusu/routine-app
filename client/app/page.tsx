import { Routine, APIResponse } from "../types/routine";
import RoutineForm from "@/components/RoutineForm";

export default async function Home() {
  const res = await fetch("http://backend:8080/api/v1/routines", {
    cache: "no-store",
  });
  const json: APIResponse<Routine[]> = await res.json();
  const routines = json.data;

  return (
    <div className="min-h-screen bg-white p-8 text-black">
      <div className="max-w-md mx-auto">
        <h1 className="text-3xl font-bold mb-6">Routine App</h1>
        <RoutineForm />

        <div className="grid gap-4">
          {routines?.map((r) => (
            <div key={r.id} className="p-4 border rounded shadow-sm bg-gray-50">
              <h2 className="text-xl font-bold">{r.title}</h2>
              <p className="text-gray-600">Interval: </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
