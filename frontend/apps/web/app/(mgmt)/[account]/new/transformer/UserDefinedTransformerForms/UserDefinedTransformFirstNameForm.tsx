'use client';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { Switch } from '@/components/ui/switch';
import { ReactElement } from 'react';
import { useFormContext } from 'react-hook-form';
import {
  CreateUserDefinedTransformerFormValues,
  UpdateUserDefinedTransformerFormValues,
} from '../schema';

interface Props {
  isDisabled?: boolean;
}

export default function UserDefinedTransformFirstNameForm(
  props: Props
): ReactElement {
  const fc = useFormContext<
    | UpdateUserDefinedTransformerFormValues
    | CreateUserDefinedTransformerFormValues
  >();
  const { isDisabled } = props;
  return (
    <div className="flex flex-col w-full space-y-4 pt-4">
      <FormField
        name={`config.value.preserveLength`}
        control={fc.control}
        render={({ field }) => (
          <FormItem className="flex flex-row items-center justify-between rounded-lg border dark:border-gray-700 p-3 shadow-sm">
            <div className="space-y-0.5">
              <FormLabel>Preserve Length</FormLabel>
              <FormDescription className="w-[90%]">
                Set the length of the output first name to be the same as the
                input
              </FormDescription>
            </div>
            <FormControl>
              <Switch
                checked={field.value}
                onCheckedChange={field.onChange}
                disabled={isDisabled}
              />
            </FormControl>
          </FormItem>
        )}
      />
    </div>
  );
}
