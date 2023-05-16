import { useVOnboarding } from "v-onboarding"
import { Ref, ref } from "vue"

export const useOnboarding = (drawer: Ref<boolean>) => {
    const wrapper = ref(null)
    const onboarding = ref(true)
    const { start, finish } = useVOnboarding(wrapper)

    const steps = [
        {
            attachTo: { element: '#example-menu' },
            content: { title: "Example Menu", description: "Open a list of predefined Policies." }
        },
        {
            attachTo: { element: '#tutorials' },
            content: { title: "Examples", description: "Choose between different Categories" },
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
            content: { title: "Policy Panel", description: "Write and modify your Kyverno Policies." }
        },
        {
            attachTo: { element: '#context-panel' },
            content: { title: "Context Panel", description: "Define or change your Context." }
        },
        {
            attachTo: { element: '#resource-panel' },
            content: { title: "Resource Panel", description: "Write and modify Resources which your Policies are applied to." }
        },
        {
            attachTo: { element: '#share-button' },
            content: { title: "Share Button", description: "Share your policies, resources and context with the community." }
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
            content: { title: "Start Button", description: "Evaluates Policies against Resources with the given Context." }
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