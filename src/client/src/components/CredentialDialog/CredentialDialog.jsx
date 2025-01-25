import { useState } from "react";
import {
  Dialog,
  DialogBackdrop,
  DialogPanel,
  DialogTitle,
} from "@headlessui/react";
import PropTypes from "prop-types";

export default function CredentialDialog({ open }) {
  const setOpen = useState(open)[1];
  return (
    <Dialog open={open} onClose={setOpen} className="relative z-10">
      <DialogBackdrop
        transition
        className="fixed inset-0 bg-black/75 backdrop-blur-xs transition-opacity data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out data-leave:duration-200 data-leave:ease-in"
      />
      <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
        <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
          <DialogPanel
            transition
            className="relative transform overflow-hidden rounded-lg text-left shadow-xl transition-all data-closed:translate-y-4 data-closed:opacity-0 data-enter:duration-300 data-enter:ease-out 
            data-leave:duration-200 data-leave:ease-in sm:my-8 sm:w-full sm:max-w-lg data-closed:sm:translate-y-0 data-closed:sm:scale-95"
          >
            <div className="px-4 pt-5 pb-4 sm:p-6 sm:pb-4 p-5 bg-gray-900 antialiased">
              <div className="sm:flex sm:items-start">
                <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                  <DialogTitle as="h2" className="font-semibold text-2xl">
                    Setup your Genote account
                  </DialogTitle>
                  <div className="mb-4 mt-4">
                    Enter your Genote account credentials and Discord webhook to
                    get started.
                  </div>
                  <div className="mt-2">
                    <div className="grid gap-6 mb-6 md:grid-cols-2">
                      <div>
                        <label
                          htmlFor="email"
                          className="block mb-2 text-sm font-medium text-white"
                        >
                          Genote identifier
                        </label>
                        <input
                          type="email"
                          id="email"
                          className="border text-sm rounded-lg block w-full p-2.5 bg-gray-700 border-gray-600 placeholder-gray-400 text-white"
                          placeholder="cip@usherbrooke.ca"
                          required
                        />
                      </div>
                      <div>
                        <label
                          htmlFor="password"
                          className="block mb-2 text-sm font-medium text-white"
                        >
                          Password
                        </label>
                        <input
                          type="password"
                          id="password"
                          className="border text-sm rounded-lg block w-full p-2.5 bg-gray-700 border-gray-600 placeholder-gray-400 text-white "
                          placeholder="•••••••••"
                          required
                        />
                      </div>
                    </div>
                    <div className="mb-6">
                      <label
                        htmlFor="discord-webhook"
                        className="block mb-2 text-sm font-medium text-white"
                      >
                        Discord Webhook
                      </label>
                      <input
                        type="text"
                        id="discord-webhook"
                        className="border text-sm rounded-lg block w-full p-2.5 bg-gray-700 border-gray-600 placeholder-gray-400 text-white focus:ring-blue-500 focus:border-blue-500"
                        placeholder="https://discord.com/api/webhooks/..."
                        required
                      />
                    </div>
                    <div className="flex justify-end">
                      <button
                        type="submit"
                        className="text-white focus:ring-4 focus:outline-none 
                         font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center bg-green-600 hover:bg-green-700 focus:ring-green-800"
                      >
                        Submit
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </DialogPanel>
        </div>
      </div>
    </Dialog>
  );
}

CredentialDialog.propTypes = {
  open: PropTypes.bool.isRequired,
};
