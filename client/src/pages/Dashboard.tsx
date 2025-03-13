import type React from "react";
import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import {
  Link2,
  Copy,
  Check,
  Trash2,
  Edit,
  Plus,
  Lock,
  Unlock,
  BarChart3,
} from "lucide-react";
import { toast } from "react-toastify";
import XCircle from "../components/XCircle";
import URLStatsCard from "../components/URLStatsCards";
import ActionButton from "../components/ActionButton";

interface URL {
  id: string;
  short_code: string;
  original_url: string;
  requires_auth: boolean;
}

const Dashboard: React.FC = () => {
  const [urls, setUrls] = useState<URL[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const { token } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    fetchUrls();
  }, []);

  const fetchUrls = async () => {
    try {
      const response = await fetch("http://localhost:8080/urls", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      // if token is expired or user is unauthorized
      if (response.status === 401) {
        navigate("/login");
        return;
      }

      const data = await response.json();
      if (!data.success) throw new Error("Failed to fetch URLs");

      setUrls(data.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(`http://localhost:8080/${text}`).then(
      () => {
        setCopiedId(id);
        toast.success("URL copied to clipboard!");
        setTimeout(() => setCopiedId(null), 2000);
      },
      () => {
        toast.error("Failed to copy URL");
      }
    );
  };

  if (loading) {
    return (
      <div className="min-h-[90vh] flex flex-col items-center justify-center bg-gradient-to-b from-white to-indigo-50 py-12 px-4">
        <div className="animate-spin rounded-full h-16 w-16 border-4 border-indigo-500 border-t-transparent"></div>
        <p className="mt-4 text-indigo-600 font-medium">Loading your URLs...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-[90vh] flex flex-col items-center justify-center bg-gradient-to-b from-white to-indigo-50 py-12 px-4">
        <div className="bg-white p-8 rounded-xl shadow-lg border border-red-100 max-w-md w-full text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-red-100 mb-4">
            <XCircle className="h-8 w-8 text-red-600" />
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">
            Error Loading URLs
          </h2>
          <p className="text-red-600 mb-6">{error}</p>
          <button
            onClick={fetchUrls}
            className="px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-lg hover:from-indigo-700 hover:to-purple-700 transition-all duration-200 shadow-md"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  // Stats for the dashboard
  const totalUrls = urls.length;
  const protectedUrls = urls.filter((url) => url.requires_auth).length;
  const publicUrls = totalUrls - protectedUrls;

  return (
    <div className="min-h-[90vh] bg-gradient-to-b from-white to-indigo-50 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        {/* Dashboard Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
            Your Dashboard
          </h1>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <URLStatsCard label="Total URLs" count={totalUrls} Icon={Link2} />
          <URLStatsCard label="Public URLs" count={publicUrls} Icon={Unlock} />
          <URLStatsCard
            label="Protected URLs"
            count={protectedUrls}
            Icon={Lock}
          />
        </div>

        {/* Action Bar */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 gap-4">
          <div className="flex items-center">
            <h2 className="text-xl font-semibold text-gray-900">
              Your Shortened URLs
            </h2>
          </div>
          <Link
            to="/create"
            className="inline-flex items-center px-4 py-2 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-lg hover:from-indigo-700 hover:to-purple-700 transition-all duration-200 shadow-md"
          >
            <Plus className="h-5 w-5 mr-2" />
            Create New URL
          </Link>
        </div>

        {urls.length === 0 ? (
          <div className="bg-white rounded-xl shadow-md p-8 text-center border border-indigo-100">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-indigo-100 mb-4">
              <Link2 className="h-8 w-8 text-indigo-600" />
            </div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">
              No URLs yet
            </h3>
            <p className="text-gray-600 mb-6 max-w-md mx-auto">
              Create your first shortened URL to start sharing links securely.
            </p>
            <Link
              to="/create"
              className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 text-white rounded-lg hover:from-indigo-700 hover:to-purple-700 transition-all duration-200 shadow-md"
            >
              <Plus className="h-5 w-5 mr-2" />
              Create Your First URL
            </Link>
          </div>
        ) : (
          <div className="bg-white rounded-xl shadow-md overflow-hidden border border-indigo-100">
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Short URL
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Original URL
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Security
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {urls.map((url) => (
                    <tr
                      key={url.id}
                      className="hover:bg-gray-50 transition-colors"
                    >
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center">
                          <Link2 className="h-4 w-4 text-indigo-500 mr-2 flex-shrink-0" />
                          <a
                            href={`http://localhost:8080/${url.short_code}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-indigo-600 hover:text-indigo-800 font-medium flex items-center"
                          >
                            {url.short_code}
                          </a>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="text-sm text-gray-900 truncate max-w-xs">
                          <a
                            href={url.original_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="hover:text-indigo-600 transition-colors"
                            title={url.original_url}
                          >
                            {url.original_url}
                          </a>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {url.requires_auth ? (
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                            <Lock className="h-3 w-3 mr-1" />
                            Protected
                          </span>
                        ) : (
                          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                            <Unlock className="h-3 w-3 mr-1" />
                            Public
                          </span>
                        )}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        <div className="flex space-x-2">
                          <ActionButton
                            onClick={() =>
                              copyToClipboard(url.short_code, url.id)
                            }
                            title="Copy URL"
                          >
                            {copiedId === url.id ? (
                              <Check className="h-5 w-5 text-green-500" />
                            ) : (
                              <Copy className="h-5 w-5" />
                            )}
                          </ActionButton>
                          {/* TODO: add the correct route for editing, deleting and viewing stats for the URL */}
                          <Link
                            to="/dashboard"
                            className="text-gray-600 hover:text-indigo-600 transition-colors p-1 rounded-md hover:bg-indigo-50"
                            title="Edit URL"
                          >
                            <Edit className="h-5 w-5" />
                          </Link>
                          <ActionButton
                            className="text-gray-600 hover:text-red-600 transition-colors p-1 rounded-md hover:bg-red-50"
                            title="Delete URL"
                          >
                            <Trash2 className="h-5 w-5" />
                          </ActionButton>
                          <Link
                            to="/dashboard"
                            className="text-gray-600 hover:text-indigo-600 transition-colors p-1 rounded-md hover:bg-indigo-50"
                            title="View Statistics"
                          >
                            <BarChart3 className="h-5 w-5" />
                          </Link>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;
