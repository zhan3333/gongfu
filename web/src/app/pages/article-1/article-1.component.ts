import { Component } from '@angular/core';
import { CommonModule, NgOptimizedImage } from '@angular/common';
import { ImageComponent } from '../image/image.component';

@Component({
  selector: 'app-article-1',
  standalone: true,
  imports: [CommonModule, NgOptimizedImage, ImageComponent],
  templateUrl: './article-1.component.html',
})
export class Article1Component {

}
