import { useCallback, useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { toast } from "react-toastify";
import { Link2, ArrowRight, Lock, Bell, Mail, Users, Info } from "lucide-react";

export default function CreateURL() {
  const [formData, setFormData] = useState({
    originalUrl: "",
    requiresAuth: false,
    notifyOnAccess: false,
    authorizedEmails: "",
  });
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();
  const { token } = useAuth();

  const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const { id, value, type, checked } = e.target;
    setFormData((prev) => ({
      ...prev,
      [id]: type === "checkbox" ? checked : value,
    }));
  }, []);

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
          original_url: formData.originalUrl,
          options: {
            requires_auth: formData.requiresAuth,
            notify_on_access: formData.notifyOnAccess,
            authorized_emails: formData.requiresAuth
              ? formData.authorizedEmails
                  .split(",")
                  .map((email) => email.trim())
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
    <div className="min-h-[90vh] flex flex-col items-center justify-center py-12 px-4 sm:px-6 lg:px-8 bg-gradient-to-b from-white to-indigo-50">
      <div className="max-w-xl w-full space-y-8">
        <div className="text-center">
          <h2 className="mt-4 text-3xl font-extrabold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
            Create Short URL
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Transform your long URLs into short, memorable links with optional
            authentication
          </p>
        </div>

        <div className="bg-white p-8 rounded-xl shadow-lg border border-indigo-100">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label
                htmlFor="originalUrl"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                Original URL
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Link2 className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  id="originalUrl"
                  type="url"
                  value={formData.originalUrl}
                  onChange={handleChange}
                  className="appearance-none block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition duration-150 ease-in-out"
                  placeholder="https://example.com/very-long-url-that-needs-shortening"
                  required
                />
              </div>
            </div>

            <div className="bg-indigo-50 p-4 rounded-lg border border-indigo-100">
              <div className="flex items-center mb-4">
                <input
                  type="checkbox"
                  id="requiresAuth"
                  checked={formData.requiresAuth}
                  onChange={handleChange}
                  className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label
                  htmlFor="requiresAuth"
                  className="ml-2 text-sm font-medium text-gray-700 flex items-center"
                >
                  <Lock className="h-4 w-4 text-indigo-600 mr-1" />
                  Require Authentication
                </label>
                <div className="ml-auto">
                  <div className="relative group">
                    <Info className="h-4 w-4 text-gray-400 cursor-help" />
                    <div className="absolute right-0 w-64 p-2 bg-white rounded shadow-lg text-xs text-gray-600 hidden group-hover:block z-10">
                      When enabled, users will need to log in before accessing
                      your link
                    </div>
                  </div>
                </div>
              </div>

              {formData.requiresAuth && (
                <div className="ml-6 mb-4">
                  <label
                    htmlFor="authorizedEmails"
                    className="text-sm font-medium text-gray-700 mb-1 flex items-center"
                  >
                    <Users className="h-4 w-4 text-indigo-600 mr-1" />
                    Authorized Emails (comma-separated)
                  </label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Mail className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      id="authorizedEmails"
                      type="text"
                      value={formData.authorizedEmails}
                      onChange={handleChange}
                      className="appearance-none block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition duration-150 ease-in-out"
                      placeholder="user@example.com, another@example.com"
                      required
                    />
                  </div>
                  <p className="mt-1 text-xs text-gray-500">
                    Leave blank to allow any authenticated user
                  </p>
                </div>
              )}

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="notifyOnAccess"
                  checked={formData.notifyOnAccess}
                  onChange={handleChange}
                  className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label
                  htmlFor="notifyOnAccess"
                  className="ml-2 text-sm font-medium text-gray-700 flex items-center"
                >
                  <Bell className="h-4 w-4 text-indigo-600 mr-1" />
                  Notify on Access
                </label>
                <div className="ml-auto">
                  <div className="relative group">
                    <Info className="h-4 w-4 text-gray-400 cursor-help" />
                    <div className="absolute right-0 w-64 p-2 bg-white rounded shadow-lg text-xs text-gray-600 hidden group-hover:block z-10">
                      You'll receive an email notification when someone accesses
                      your link
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div className="pt-4">
              <button
                type="submit"
                disabled={loading}
                className="group relative w-full flex justify-center py-3 px-4 border border-transparent rounded-lg text-white bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-150 ease-in-out font-medium shadow-sm disabled:opacity-70 cursor-pointer"
              >
                {loading ? (
                  <span className="flex items-center">
                    <svg
                      className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      ></circle>
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                    Creating URL...
                  </span>
                ) : (
                  <span className="flex items-center">
                    Create Short URL
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </span>
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
