import { useVOnboarding } from 'v-onboarding'
import { Ref, ref } from 'vue'

export const useOnboarding = (drawer: Ref<boolean>, advanced: Ref<boolean>) => {
  const wrapper = ref(null)
  const onboarding = ref(true)
  const { start, finish } = useVOnboarding(wrapper)

  const steps = [
    {
      attachTo: { element: '#example-menu' },
      content: {
        title: 'Example Menu',
        description:
          'The Examples menu contains a list of sample Kyverno policies. Selecting a policy will populate not only the policy itself but also test resources and, in some cases, a runtime context for convenient and quick testing.'
      }
    },
    {
      attachTo: { element: '#tutorials' },
      content: {
        title: 'Examples',
        description:
          'Choose between different categories of policies. The Tutorials category contains several Kyverno policies grouped by type which demonstrate how to use the Playground with various different policy styles.'
      },
      on: {
        beforeStep: async () => {
          if (drawer.value) return

          drawer.value = true
          await new Promise((r) => setTimeout(r, 200))
        },
        afterStep: () => {
          drawer.value = false
        }
      }
    },
    {
      attachTo: { element: '#advanced-btn' },
      content: {
        title: 'Advanced Configuration',
        description: 'The advanced configuration is used to simulate more complex policies.'
      },
      on: {
        beforeStep: async () => {
          if (advanced.value) return

          advanced.value = true
          await new Promise((r) => setTimeout(r, 200))
        }
      }
    },
    {
      attachTo: { element: '#context-btn' },
      content: {
        title: 'Context Configuration',
        description:
          'Contexts are special metadata used to define the runtime context of the test and include things like the Kubernetes version, the metadata of the AdmissionReview request outside of the resource itself (if required), and any variables which may need to be statically defined. Variables which begin with request.object do not need to be defined here.'
      }
    },
    {
      attachTo: { element: '#config-btn' },
      content: {
        title: 'Kyverno Configuration',
        description: 'The Kyverno configuration can be changed to update resource filters or the default registry.'
      }
    },
    {
      attachTo: { element: '#crd-btn' },
      content: {
        title: 'Custom Resource Definitions',
        description: 'Add CRDs when testing custom resources.'
      }
    },
    {
      attachTo: { element: '#exceptions-btn' },
      content: {
        title: 'Policy Exceptions',
        description: 'Test Policy Exceptions by providing them in this panel.'
      }
    },
    {
      attachTo: { element: '#clusterResources-btn' },
      content: {
        title: 'Cluster Resources',
        description: 'Simulate existing resources in the cluster to test variable substitution or when testing generate rules using the clone declaration.'
      },
      on: {
        afterStep: () => {
          advanced.value = false
        }
      }
    },
    {
      attachTo: { element: '#policy-panel' },
      content: {
        title: 'Policy Panel',
        description:
          'Kyverno policies are placed in the Policy Panel where they can be modified and tested. Schema validation of the editor updates in real time to show you any errors that may be found in the policy.'
      }
    },
    {
      attachTo: { element: '#resource-panel' },
      content: {
        title: 'Resource Panel',
        description:
          "Resources are where you define the Kubernetes resources which are tested against the policy defined in the Policies pane. Multiple resources are supported with the standard YAML document delimiter '---'. A split pane button enables you to provide old and new versions of a given resource when testing UPDATE requests."
      }
    },
    {
      attachTo: { element: '#share-button' },
      content: {
        title: 'Share Button',
        description:
          'Share your policies, resources, and context with the community. A link will be produced which fully encodes all the panels and their contents making it simple to show others everything you have built in the Playground.'
      }
    },
    {
      attachTo: { element: '#save-button' },
      content: { title: 'Save Button', description: 'Save your tests locally as named profiles for further testing at a later time.' }
    },
    {
      attachTo: { element: '#load-button' },
      content: { title: 'Load Button', description: 'Load your local persisted profiles or reset your inputs with the default profile.' }
    },
    {
      attachTo: { element: '#start-btn' },
      content: {
        title: 'Start Button',
        description: 'Begin testing policies against resources with the provided context. The Results window will return the results of all provided resources.'
      }
    }
  ]

  return {
    finish,
    wrapper,
    start,
    onboarding,
    steps
  }
}
