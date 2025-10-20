import { ErrorInput } from "@/types/response";

export default function ErrorInputMessage({
  errors,
  field,
}: {
  errors?: ErrorInput[];
  field: string;
}) {
  return (
    <div className="mx-4">
      {errors?.map((error) =>
        error.field === field ? (
          <ul key={error.field} className="list-disc text-xs text-red-500">
            {error.reasons.map((reason) => (
              <li key={reason}>{reason}</li>
            ))}
          </ul>
        ) : null,
      )}
    </div>
  );
}
