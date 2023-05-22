import { useVOnboarding } from "v-onboarding"
import { Ref, ref } from "vue"

export const useOnboarding = (drawer: Ref<boolean>) => {
    const wrapper = ref(null)
    const onboarding = ref(true)
    const { start, finish } = useVOnboarding(wrapper)

    const steps = [
        {
            attachTo: { element: '#example-menu' },
            content: { title: "Example Menu", description: "The Examples menu contains a list of predefined Kyverno policies with test resources and contexts." }
        },
        {
            attachTo: { element: '#tutorials' },
            content: { title: "Examples", description: "Choose between different Categories. The Tutorials category contains several Kyverno policies and supporting resources which demonstrate how to use the Playground effectively." },
            on: {
                beforeStep: async () => {
                    drawer.value = true;
                    await new Promise(r => setTimeout(r, 200))
                },
                afterStep: () => {
                    drawer.value = false
                }
            }
        },
        {
            attachTo: { element: '#policy-panel' },
            content: { title: "Policy Panel", description: "A Kyverno policy goes in the Policy Panel where it can be modified and tested. Schema validation of the editor updates in real time to show you any errors that may be found in the policy." }
        },
        {
            attachTo: { element: '#context-panel' },
            content: { title: "Context Panel", description: "Contexts are special metadata used to define the runtime context of the test and include things like the Kubernetes version, the metadata of the AdmissionReview request outside of the resource itself (if required), and any variables which may need to be statically defined. Variables which begin with request.object do not need to be defined here. You can collapse the Context panel to save screen space if you wish." }
        },
        {
            attachTo: { element: '#resource-panel' },
            content: { title: "Resource Panel", description: "Resources are where you define the Kubernetes resources which are tested against the policy defined in the Policies pane. Multiple resources are supported with the standard YAML document delimiter '---'." }
        },
        {
            attachTo: { element: '#share-button' },
            content: { title: "Share Button", description: "Share your policies, resources, and context with the community. A link will be produced which fully encodes all the panels and their contents making it simple to show others everything you have built in the Playground." }
        },
        {
            attachTo: { element: '#save-button' },
            content: { title: "Save Button", description: "Save your tests locally as named profiles for further testing at a later time." }
        },
        {
            attachTo: { element: '#load-button' },
            content: { title: "Load Button", description: "Load your local persisted profiles or reset your inputs with the default profile." }
        },
        {
            attachTo: { element: '#start-btn' },
            content: { title: "Start Button", description: "Begin testing policies against resources with the provided context. The Results window will return the results of all provided resources." }
        },
    ]

    return {
        finish,
        wrapper,
        start,
        onboarding,
        steps
    }
}