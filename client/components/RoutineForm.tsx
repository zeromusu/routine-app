"use client";

import { useRouter } from "next/navigation";
import React, { useState } from "react";

export default function RoutineForm() {
  const router = useRouter();
  const [title, setTitle] = useState("");
  const [interval, setInterval] = useState("daily");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const res = await fetch("http://localhost:8080/api/v1/routines", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ title, interval }),
    });

    if (res.ok) {
      setTitle("");
      router.refresh();
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mb-8 p-6 bg-gray-100 rounded-lg">
      <div className="flex flex-col gap-4">
        <input
          type="text"
          placeholder="Input Routine"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="p-2 border rounded text-black"
          required
        />
        <select
          value={interval}
          onChange={(e) => setInterval(e.target.value)}
          className="p-2 border rounded text-black"
        >
          <option value="daily">Daily</option>
          <option value="weekly">Weekly</option>
          <option value="monthly">Monthly</option>
        </select>
        <button
          type="submit"
          className="bg-blue-600 text-white p-2 rounded hover:bg-blue-700"
        >
          Add Routine
        </button>
      </div>
    </form>
  );
}
