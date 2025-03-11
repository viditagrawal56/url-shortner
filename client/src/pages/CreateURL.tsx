import { useState, FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { toast } from "react-toastify";

export default function CreateURL() {
  const [originalUrl, setOriginalUrl] = useState("");
  const [requiresAuth, setRequiresAuth] = useState(false);
  const [notifyOnAccess, setNotifyOnAccess] = useState(false);
  const [authorizedEmails, setAuthorizedEmails] = useState("");
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();
  const { token } = useAuth();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await fetch("http://localhost:8080/urls", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          original_url: originalUrl,
          options: {
            requires_auth: requiresAuth,
            notify_on_access: notifyOnAccess,
            authorized_emails: requiresAuth
              ? authorizedEmails.split(",").map((email) => email.trim())
              : [],
          },
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Failed to create short URL");
      }

      toast.success("URL shortened successfully!");
      navigate("/dashboard");
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to create short URL"
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-lg mx-auto mt-10 p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-6 text-center">Create Short URL</h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700">
            Original URL
          </label>
          <input
            type="url"
            value={originalUrl}
            onChange={(e) => setOriginalUrl(e.target.value)}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            required
          />
        </div>

        <div className="flex items-center">
          <input
            type="checkbox"
            id="requiresAuth"
            checked={requiresAuth}
            onChange={(e) => setRequiresAuth(e.target.checked)}
            className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
          />
          <label htmlFor="requiresAuth" className="ml-2 text-sm text-gray-700">
            Require Authentication
          </label>
        </div>

        {requiresAuth && (
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Authorized Emails (comma-separated)
            </label>
            <input
              type="text"
              value={authorizedEmails}
              onChange={(e) => setAuthorizedEmails(e.target.value)}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
          </div>
        )}

        <div className="flex items-center">
          <input
            type="checkbox"
            id="notifyOnAccess"
            checked={notifyOnAccess}
            onChange={(e) => setNotifyOnAccess(e.target.checked)}
            className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
          />
          <label
            htmlFor="notifyOnAccess"
            className="ml-2 text-sm text-gray-700"
          >
            Notify on Access
          </label>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
        >
          {loading ? "Creating..." : "Create Short URL"}
        </button>
      </form>
    </div>
  );
}
