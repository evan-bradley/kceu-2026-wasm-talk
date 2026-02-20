import { ExportResult, ExportResultCode } from '@opentelemetry/core';
import {
  AggregationTemporality,
  InstrumentType,
  PushMetricExporter,
  ResourceMetrics,
} from '@opentelemetry/sdk-metrics';
import { JsonMetricsSerializer } from '@opentelemetry/otlp-transformer';

/**
 * A custom PushMetricExporter that serializes SDK metrics to OTLP JSON
 * and sends them to the OTel Collector running in WebAssembly via
 * globalThis.__otelReceiveCallback.
 */
export class WasmCollectorExporter implements PushMetricExporter {
  private _shutdown = false;

  export(
    metrics: ResourceMetrics,
    resultCallback: (result: ExportResult) => void
  ): void {
    if (this._shutdown) {
      resultCallback({ code: ExportResultCode.FAILED });
      return;
    }

    try {
      // Serialize SDK ResourceMetrics to OTLP JSON (Uint8Array)
      const serialized = JsonMetricsSerializer.serializeRequest(metrics);

      // Send to the Collector's JS receiver via js.CopyBytesToGo
      if ((globalThis as any).__otelReceiveCallback) {
        (globalThis as any).__otelReceiveCallback(serialized);
      }

      resultCallback({ code: ExportResultCode.SUCCESS });
    } catch (e) {
      resultCallback({
        code: ExportResultCode.FAILED,
        error: e instanceof Error ? e : new Error(String(e)),
      });
    }
  }

  forceFlush(): Promise<void> {
    return Promise.resolve();
  }

  selectAggregationTemporality(
    _instrumentType: InstrumentType
  ): AggregationTemporality {
    return AggregationTemporality.DELTA;
  }

  shutdown(): Promise<void> {
    this._shutdown = true;
    return Promise.resolve();
  }
}
