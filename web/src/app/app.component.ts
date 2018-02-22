import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

const API = '/api';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  seeds = [];
  remove = [];
  seedModal;
  seed = {};

  constructor(public http: HttpClient) {
    this.getSeeds();
  }

  getSeeds = () => {
    this.http.get(API + '/seeds').subscribe(res => this.seeds = <any>res);
  }

  showSeed = (seed: string) => {
    this.http.get(API + '/seed', { params: {seed: seed} }).subscribe(res => {
      if (res['seed_id']) {
        this.seed = res;
        this.seedModal = true;
      } else {
        this.seed = {};
        this.getSeeds();
        console.log('Seed not found');
      }
    });
  }

  save = (seedId, seedSecret) => {
    if (seedId.value && seedSecret.value) {
      this.http.put(API + '/seed', {seed_id: seedId.value, seed_secret: seedSecret.value}).subscribe(res => {
        this.seed = {};
        this.seedModal = false;
        this.getSeeds();
      }, err => {});
    }
  }

  dismiss = () => {
    this.seed = {};
    this.seedModal = false;
  }

  createSeed = () => {
    this.seed = {seed_id: '', seed_secret: ''};
    this.seedModal = true;
  }

  cancelRemoval = (seed: string) => {
    this.remove = this.remove.filter(s => s !== seed);
  }

  removeSeed = (seed: string, confirm?: boolean) => {
    if (!confirm) {
      this.remove.push(seed);
    } else {
      this.http.delete(API + '/remove', { params: { seed: seed } }).subscribe(res => {
        this.cancelRemoval(seed);
        this.getSeeds();
      }, err => {});
    }
  }
}
