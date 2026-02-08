import { Routine, APIResponse } from "../types/routine";
import RoutineForm from "@/components/RoutineForm";
import EditableRoutine from "@/components/EditableRoutine";

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
            <EditableRoutine key={r.id} routine={r} />
          ))}
        </div>
      </div>
    </div>
  );
}
