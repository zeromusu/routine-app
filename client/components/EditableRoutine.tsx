"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { Routine } from "@/types/routine";
import DeleteButton from "./DeleteButton";

export default function EditableRoutine({ routine }: { routine: Routine }) {
  const router = useRouter();
  const [isEditing, setIsEditing] = useState(false);
  const [title, setTitle] = useState(routine.title);
  const [interval, setInterval] = useState(routine.interval);

  const handleUpdate = async () => {
    const res = await fetch(
      `http://localhost:8080/api/v1/routines/${routine.id}`,
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, interval }),
      },
    );

    if (res.ok) {
      setIsEditing(false);
      router.refresh();
    }
  };

  if (isEditing) {
    return (
      <div className="p-4 border rounded shadow-md bg-blue-50 flex flex-col gap-2">
        <input
          className="border p-1 rounded text-black"
          onChange={(e) => setTitle(e.target.value)}
        />
        <select
          className="border p-1 rounded text-black"
          value={interval}
          onChange={(e) => setInterval(e.target.value)}
        >
          <option value="daily">daily</option>
          <option value="weekly">weekly</option>
          <option value="monthly">monthly</option>
        </select>
        <div className="flex gap-2 justify-end">
          <button
            onClick={() => setIsEditing(false)}
            className="text-gray-500 text-sm"
          >
            cancel
          </button>
          <button
            onClick={handleUpdate}
            className="bg-blue-500 text-white px-3 py-1 rounded text-sm"
          >
            save
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-4 border rounded shadow-sm bg-gray-50 flex justify-between items-center">
      <div>
        <h2 className="text-xl font-bold">{routine.title}</h2>
        <p className="text-gray-600 text-sm">{routine.interval}</p>
      </div>
      <div className="flex gap-2">
        <button onClick={() => setIsEditing(true)} className="text-blue-500">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
            />
          </svg>
        </button>
        <DeleteButton id={routine.id} />
      </div>
    </div>
  );
}
