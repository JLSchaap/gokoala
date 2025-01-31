import { CommonModule } from '@angular/common'
import { Component, ElementRef, Input, OnChanges, OnInit, ViewEncapsulation } from '@angular/core'
import { recordStyleLayer } from 'ol-mapbox-style'
import { NgChanges } from '../vectortile-view/vectortile-view.component'
import { LegendItemComponent } from './legend-item/legend-item.component'
import { LegendItem, MapboxStyle, MapboxStyleService } from '../mapbox-style.service'
import { NGXLogger } from 'ngx-logger'

@Component({
  selector: 'app-legend-view',
  templateUrl: './legend-view.component.html',
  styleUrls: ['./legend-view.component.css'],
  imports: [CommonModule, LegendItemComponent],
  standalone: true,
  encapsulation: ViewEncapsulation.Emulated,
})
export class LegendViewComponent implements OnInit, OnChanges {
  mapboxStyle!: MapboxStyle
  @Input() styleUrl!: string
  @Input() titleItems!: string //= 'type,plus_type,functie,fysiek_voorkomen'

  LegendItems: LegendItem[] = []

  constructor(
    private logger: NGXLogger,
    private mapboxStyleService: MapboxStyleService,
    private elementRef: ElementRef
  ) {
    recordStyleLayer(true)
  }

  ngOnChanges(changes: NgChanges<LegendViewComponent>) {
    if (changes.styleUrl?.previousValue !== changes.styleUrl?.currentValue) {
      if (!changes.styleUrl.isFirstChange()) {
        this.generateLegend()
      }
    }
    if (this.titleItems) {
      if (changes.titleItems.previousValue !== changes.titleItems.currentValue) {
        if (!changes.titleItems.isFirstChange()) {
          this.generateLegend()
        }
      }
    }
  }

  ngOnInit(): void {
    this.generateLegend()
  }

  generateLegend() {
    if (this.styleUrl) {
      this.mapboxStyleService.getMapboxStyle(this.styleUrl).subscribe(style => {
        this.mapboxStyle = this.mapboxStyleService.removeRasterLayers(style)
        if (this.titleItems) {
          const titlepart = this.titleItems.split(',')
          this.LegendItems = this.mapboxStyleService.getItems(this.mapboxStyle, this.mapboxStyleService.customTitle, titlepart)
        } else {
          this.LegendItems = this.mapboxStyleService.getItems(this.mapboxStyle, this.mapboxStyleService.capitalizeFirstLetter, [])
        }
      })
    } else {
      this.logger.error('no style url supplied')
    }
  }
}
