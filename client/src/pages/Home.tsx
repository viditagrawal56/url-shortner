import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import {
  ArrowRight,
  Lock,
  Link,
  Share2,
  Users,
  ChevronRight,
} from "lucide-react";

function Home() {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const handleButtonClick = () => {
    if (isAuthenticated) {
      navigate("/dashboard");
    } else {
      navigate("/login");
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-indigo-50 via-white to-purple-50">
      {/* Hero Section */}
      <section className="pt-20 pb-16 md:pt-28 md:pb-24">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row items-center">
            <div className="md:w-1/2 md:pr-12">
              <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-gray-900 leading-tight">
                Short Links with{" "}
                <span className="text-indigo-600">Private Access</span>
              </h1>
              <p className="mt-6 text-lg text-gray-600 leading-relaxed">
                Create short, memorable links that can be protected with
                authentication. Perfect for sharing sensitive content with only
                the people you choose.
              </p>
              <div className="mt-8 flex flex-col sm:flex-row gap-4">
                <button
                  onClick={handleButtonClick}
                  className="px-8 py-4 bg-indigo-600 text-white font-semibold rounded-lg shadow-lg hover:bg-indigo-700 transition duration-300 flex items-center justify-center"
                >
                  {isAuthenticated ? "Go to Dashboard" : "Get Started"}
                  <ArrowRight className="ml-2 h-5 w-5" />
                </button>
                <button
                  onClick={() => navigate("/login")}
                  className="px-8 py-4 bg-white text-indigo-600 font-semibold rounded-lg shadow-md border border-indigo-100 hover:bg-indigo-50 transition duration-300"
                >
                  Learn More
                </button>
              </div>
            </div>
            <div className="md:w-1/2 mt-12 md:mt-0">
              <div className="relative">
                <div className="absolute inset-0 bg-gradient-to-r from-indigo-500 to-purple-500 rounded-3xl transform rotate-3 scale-105 opacity-10"></div>
                <div className="relative bg-white rounded-2xl shadow-xl p-6 md:p-8">
                  <div className="flex items-center mb-6">
                    <div className="w-3 h-3 bg-red-500 rounded-full mr-2"></div>
                    <div className="w-3 h-3 bg-yellow-500 rounded-full mr-2"></div>
                    <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                    <div className="ml-4 bg-gray-100 rounded-md px-4 py-2 flex-1 text-sm text-gray-500">
                      https://securelink.io
                    </div>
                  </div>
                  <div className="space-y-4">
                    <div className="bg-gray-50 p-4 rounded-lg border border-gray-100">
                      <h3 className="font-medium text-gray-800">
                        Create Secure Link
                      </h3>
                      <div className="mt-3 flex">
                        <input
                          type="text"
                          placeholder="Paste your long URL here"
                          className="flex-1 p-3 border border-gray-200 rounded-l-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 text-sm"
                        />
                        <button className="bg-indigo-600 text-white px-4 rounded-r-lg hover:bg-indigo-700 transition">
                          Shorten
                        </button>
                      </div>
                      <div className="mt-3 flex items-center">
                        <input
                          type="checkbox"
                          id="require-auth"
                          className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                        />
                        <label
                          htmlFor="require-auth"
                          className="ml-2 text-sm text-gray-600"
                        >
                          Require authentication to access
                        </label>
                      </div>
                    </div>
                    <div className="bg-indigo-50 p-4 rounded-lg border border-indigo-100 flex items-center">
                      <Lock className="h-5 w-5 text-indigo-600 mr-3" />
                      <div>
                        <p className="text-sm text-indigo-900">
                          Your secure link will only be accessible to authorized
                          users
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 bg-white">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900">
              Powerful Features for Secure Sharing
            </h2>
            <p className="mt-4 text-lg text-gray-600 max-w-2xl mx-auto">
              Our platform offers unique capabilities that make sharing links
              both convenient and secure.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="bg-gradient-to-br from-white to-indigo-50 p-8 rounded-xl shadow-md hover:shadow-lg transition duration-300">
              <div className="bg-indigo-100 p-3 rounded-full w-14 h-14 flex items-center justify-center mb-6">
                <Lock className="h-7 w-7 text-indigo-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">
                Authentication Protected
              </h3>
              <p className="text-gray-600">
                Create links that require authentication before access, perfect
                for sharing sensitive content.
              </p>
            </div>

            <div className="bg-gradient-to-br from-white to-indigo-50 p-8 rounded-xl shadow-md hover:shadow-lg transition duration-300">
              <div className="bg-indigo-100 p-3 rounded-full w-14 h-14 flex items-center justify-center mb-6">
                <Link className="h-7 w-7 text-indigo-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">
                Link Open Notifications
              </h3>
              <p className="text-gray-600">
                Get notified via email whenever someone opens your shared short
                link.
              </p>
            </div>

            <div className="bg-gradient-to-br from-white to-indigo-50 p-8 rounded-xl shadow-md hover:shadow-lg transition duration-300">
              <div className="bg-indigo-100 p-3 rounded-full w-14 h-14 flex items-center justify-center mb-6">
                <Share2 className="h-7 w-7 text-indigo-600" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-3">
                Analytics & Tracking
              </h3>
              <p className="text-gray-600">
                Track who's accessing your links with detailed analytics and
                insights.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="py-16 bg-gradient-to-b from-white to-indigo-50">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900">
              How It Works
            </h2>
            <p className="mt-4 text-lg text-gray-600 max-w-2xl mx-auto">
              Creating and sharing secure links is simple with our platform.
            </p>
          </div>

          <div className="max-w-4xl mx-auto">
            <div className="flex flex-col md:flex-row items-center mb-12">
              <div className="md:w-1/3 mb-6 md:mb-0">
                <div className="bg-white w-16 h-16 rounded-full shadow-md flex items-center justify-center mx-auto border-2 border-indigo-100">
                  <span className="text-2xl font-bold text-indigo-600">1</span>
                </div>
              </div>
              <div className="md:w-2/3 md:pl-8">
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  Paste Your Long URL
                </h3>
                <p className="text-gray-600">
                  Start by pasting the long URL you want to shorten into our
                  easy-to-use platform.
                </p>
              </div>
            </div>

            <div className="flex flex-col md:flex-row items-center mb-12">
              <div className="md:w-1/3 mb-6 md:mb-0">
                <div className="bg-white w-16 h-16 rounded-full shadow-md flex items-center justify-center mx-auto border-2 border-indigo-100">
                  <span className="text-2xl font-bold text-indigo-600">2</span>
                </div>
              </div>
              <div className="md:w-2/3 md:pl-8">
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  Choose Security Level
                </h3>
                <p className="text-gray-600">
                  Select whether you want your link to be public or require
                  authentication for access.
                </p>
              </div>
            </div>

            <div className="flex flex-col md:flex-row items-center">
              <div className="md:w-1/3 mb-6 md:mb-0">
                <div className="bg-white w-16 h-16 rounded-full shadow-md flex items-center justify-center mx-auto border-2 border-indigo-100">
                  <span className="text-2xl font-bold text-indigo-600">3</span>
                </div>
              </div>
              <div className="md:w-2/3 md:pl-8">
                <h3 className="text-xl font-semibold text-gray-900 mb-2">
                  Share With Confidence
                </h3>
                <p className="text-gray-600">
                  Share your secure short link with friends, family, or
                  colleagues knowing only authorized users can access it.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Testimonials */}
      <section className="py-16 bg-white">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900">
              What Our Users Say
            </h2>
            <p className="mt-4 text-lg text-gray-600 max-w-2xl mx-auto">
              Join thousands of satisfied users who trust our platform for
              secure link sharing.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="bg-white p-8 rounded-xl shadow-md border border-gray-100">
              <div className="flex items-center mb-4">
                <div className="w-12 h-12 bg-indigo-100 rounded-full flex items-center justify-center">
                  <Users className="h-6 w-6 text-indigo-600" />
                </div>
                <div className="ml-4">
                  <h4 className="font-semibold text-gray-900">Sarah Johnson</h4>
                  <p className="text-sm text-gray-500">Marketing Director</p>
                </div>
              </div>
              <p className="text-gray-600 italic">
                "The authentication feature is a game-changer. I can now share
                confidential marketing materials with my team using short links
                that only they can access."
              </p>
            </div>

            <div className="bg-white p-8 rounded-xl shadow-md border border-gray-100">
              <div className="flex items-center mb-4">
                <div className="w-12 h-12 bg-indigo-100 rounded-full flex items-center justify-center">
                  <Users className="h-6 w-6 text-indigo-600" />
                </div>
                <div className="ml-4">
                  <h4 className="font-semibold text-gray-900">Michael Chen</h4>
                  <p className="text-sm text-gray-500">Software Developer</p>
                </div>
              </div>
              <p className="text-gray-600 italic">
                "I use this service to share development resources with specific
                team members. The authentication layer adds an extra level of
                security I couldn't find elsewhere."
              </p>
            </div>

            <div className="bg-white p-8 rounded-xl shadow-md border border-gray-100">
              <div className="flex items-center mb-4">
                <div className="w-12 h-12 bg-indigo-100 rounded-full flex items-center justify-center">
                  <Users className="h-6 w-6 text-indigo-600" />
                </div>
                <div className="ml-4">
                  <h4 className="font-semibold text-gray-900">
                    Emily Rodriguez
                  </h4>
                  <p className="text-sm text-gray-500">HR Manager</p>
                </div>
              </div>
              <p className="text-gray-600 italic">
                "Perfect for sharing confidential HR documents. I can create
                short links that require authentication, ensuring sensitive
                information stays protected."
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-16 bg-gradient-to-r from-indigo-600 to-purple-600 text-white">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto text-center">
            <h2 className="text-3xl md:text-4xl font-bold mb-6">
              Ready to Create Secure Short Links?
            </h2>
            <p className="text-lg text-indigo-100 mb-8">
              Join thousands of users who trust our platform for secure,
              authenticated link sharing.
            </p>
            <button
              onClick={handleButtonClick}
              className="px-8 py-4 bg-white text-indigo-600 font-semibold rounded-lg shadow-lg hover:bg-indigo-50 transition duration-300 flex items-center mx-auto"
            >
              {isAuthenticated ? "Go to Dashboard" : "Get Started Now"}
              <ChevronRight className="ml-2 h-5 w-5" />
            </button>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-400 flex items-center justify-center">
        <div className="border-t border-gray-800 py-8 text-center text-sm">
          <p>
            &copy; {new Date().getFullYear()} SecureLink. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}

export default Home;
